package xelon

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/ssh"
	"github.com/docker/machine/libmachine/state"

	"github.com/Xelon-AG/docker-machine-driver-xelon/api"
)

const (
	defaultCPUCores       = 2
	defaultDevicePassword = "Xelon22"
	defaultDiskSize       = 20
	defaultKubernetesID   = "kub1"
	defaultMemory         = 2
	defaultSSHPort        = 22
	defaultSSHUser        = "root"
	defaultSwapDiskSize   = 2
)

type Driver struct {
	*drivers.BaseDriver
	APIBaseURL     string
	CPUCores       int
	DevicePassword string
	DiskSize       int
	KubernetesID   string
	LocalVMID      string
	Memory         int
	SwapDiskSize   int
	TenantID       string
	Token          string
}

func NewDriver(hostName, storePath string) *Driver {
	return &Driver{
		BaseDriver: &drivers.BaseDriver{
			MachineName: hostName,
			StorePath:   storePath,
		},
	}
}

func (d *Driver) Create() error {
	log.Info("Authenticating into Xelon VDC...")
	client := d.getClient()
	tenant, _, err := client.Tenant.Get()
	if err != nil {
		return err
	}
	log.Debugf("User tenant id: %v", tenant.TenantIdentifier)
	d.TenantID = tenant.TenantIdentifier

	log.Debug("(workaround): generate random delay before creating Xelon device...")
	randomDelay()

	log.Info("Creating Xelon device...")
	deviceCreateResponse, err := d.createDevice()
	if err != nil {
		return err
	}
	log.Debugf("DeviceCreateResponse: %+v", deviceCreateResponse)

	d.LocalVMID = deviceCreateResponse.Device.LocalVMID
	d.IPAddress = deviceCreateResponse.IPs[0]

	log.Info("Waiting until Xelon device will be provisioned...")
	retryCount := 5
	currentRetry := 1
	for {
		deviceRoot, _, err := client.Devices.Get(tenant.TenantIdentifier, deviceCreateResponse.Device.LocalVMID)
		if err != nil {
			log.Debugf("Error by getting device information: retry %v of %v", currentRetry, retryCount)
			if currentRetry <= retryCount {
				currentRetry++
				log.Debug("Waiting 5 seconds before next call...")
				time.Sleep(5 * time.Second)
				continue
			}
			log.Info("Xelon device could not created, clean up resources...")
			_ = d.Remove()
			return err
		}
		device := deviceRoot.Device
		toolsStatus := deviceRoot.ToolsStatus
		log.Debugf("device.powerstate: %v, device.state: %v, tools.runningStatus: %v", device.Powerstate, device.LocalVMDetails.State, toolsStatus.RunningStatus)
		if device.Powerstate == true && device.LocalVMDetails.State == 1 && toolsStatus.RunningStatus == "guestToolsRunning" {
			break
		}
		time.Sleep(2 * time.Second)
	}

	log.Debug("(workaround): waiting 15 seconds to be sure that server is ready...")
	time.Sleep(15 * time.Second)

	log.Info("Adding SSH key to the device...")
	err = d.addSSHKey(d.LocalVMID)
	if err != nil {
		return err
	}

	log.Info("Starting Xelon device...")
	err = d.startDevice()
	if err != nil {
		return err
	}

	log.Debugf("Created device LocalVMID %v, IP address %v", d.LocalVMID, d.IPAddress)

	return nil
}

func (d *Driver) DriverName() string {
	return "xelon"
}

func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	return []mcnflag.Flag{
		mcnflag.StringFlag{
			EnvVar: "XELON_API_BASE_URL",
			Name:   "xelon-api-base-url",
			Usage:  "Xelon API base URL",
		},
		mcnflag.IntFlag{
			EnvVar: "XELON_CPU_CORES",
			Name:   "xelon-cpu-cores",
			Usage:  "Number of CPU cores for the device",
			Value:  defaultCPUCores,
		},
		mcnflag.StringFlag{
			EnvVar: "XELON_DEVICE_PASSWORD",
			Name:   "xelon-device-password",
			Usage:  "Password for the device",
			Value:  defaultDevicePassword,
		},
		mcnflag.IntFlag{
			EnvVar: "XELON_DISK_SIZE",
			Name:   "xelon-disk-size",
			Usage:  "Drive size for the device in GB",
			Value:  defaultDiskSize,
		},
		mcnflag.StringFlag{
			EnvVar: "XELON_KUBERNETES_ID",
			Name:   "xelon-kubernetes-id",
			Usage:  "Kubernetes ID for the device",
			Value:  defaultKubernetesID,
		},
		mcnflag.IntFlag{
			EnvVar: "XELON_MEMORY",
			Name:   "xelon-memory",
			Usage:  "Size of memory for the device in GB",
			Value:  defaultMemory,
		},
		mcnflag.IntFlag{
			EnvVar: "XELON_SSH_PORT",
			Name:   "xelon-ssh-port",
			Usage:  "SSH port to connect",
			Value:  defaultSSHPort,
		},
		mcnflag.StringFlag{
			EnvVar: "XELON_SSH_USER",
			Name:   "xelon-ssh-user",
			Usage:  "SSH username to connect",
			Value:  defaultSSHUser,
		},
		mcnflag.IntFlag{
			EnvVar: "XELON_SWAP_DISK_SIZE",
			Name:   "xelon-swap-disk-size",
			Usage:  "Swap disk size for the device in GB",
			Value:  defaultSwapDiskSize,
		},
		mcnflag.StringFlag{
			EnvVar: "XELON_TOKEN",
			Name:   "xelon-token",
			Usage:  "Xelon authentication token",
		},
	}
}

func (d *Driver) GetSSHHostname() (string, error) {
	return d.GetIP()
}

func (d *Driver) GetURL() (string, error) {
	if err := drivers.MustBeRunning(d); err != nil {
		return "", err
	}

	ip, err := d.GetIP()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("tcp://%s", net.JoinHostPort(ip, "2376")), nil
}

func (d *Driver) GetState() (state.State, error) {
	deviceRoot, _, err := d.getClient().Devices.Get(d.TenantID, d.LocalVMID)
	if err != nil {
		return state.Error, err
	}

	if deviceRoot == nil {
		return state.None, nil
	}

	device := deviceRoot.Device
	toolsStatus := deviceRoot.ToolsStatus
	if device.Powerstate == false {
		return state.Stopped, nil
	} else {
		if device.LocalVMDetails.State == 1 && toolsStatus.RunningStatus == "guestToolsRunning" {
			return state.Running, nil
		} else {
			return state.Starting, nil
		}
	}
}

func (d *Driver) Kill() error {
	_, err := d.getClient().Devices.Stop(d.LocalVMID)
	return err
}

func (d *Driver) PreCreateCheck() error {
	if len(d.DevicePassword) < 6 {
		return fmt.Errorf("xelon-device-password must be at least 6 characters long")
	}
	return nil
}

func (d *Driver) Remove() error {
	log.Info("Stopping Xelon device...")
	err := d.stopDevice()
	if err != nil {
		return err
	}

	log.Info("Deleting Xelon device...")
	client := d.getClient()
	if resp, err := client.Devices.Delete(d.LocalVMID); err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			log.Info("Xelon device doesn't exist, assuming it is already deleted")
		} else {
			return err
		}
	}

	return nil
}

func (d *Driver) Restart() error {
	err := d.stopDevice()
	if err != nil {
		return err
	}
	return d.startDevice()
}

func (d *Driver) SetConfigFromFlags(opts drivers.DriverOptions) error {
	d.APIBaseURL = opts.String("xelon-api-base-url")
	d.CPUCores = opts.Int("xelon-cpu-cores")
	d.DevicePassword = opts.String("xelon-device-password")
	d.DiskSize = opts.Int("xelon-disk-size")
	d.KubernetesID = opts.String("xelon-kubernetes-id")
	d.Memory = opts.Int("xelon-memory")
	d.SSHPort = opts.Int("xelon-ssh-port")
	d.SSHUser = opts.String("xelon-ssh-user")
	d.SwapDiskSize = opts.Int("xelon-swap-disk-size")
	d.Token = opts.String("xelon-token")

	if d.Token == "" {
		return fmt.Errorf("xelon driver requires the --xelon-token option")
	}

	return nil
}

func (d *Driver) Start() error {
	return d.startDevice()
}

func (d *Driver) Stop() error {
	return d.stopDevice()
}

func (d *Driver) getClient() *api.Client {
	client := api.NewClient(d.Token)
	if d.APIBaseURL != "" {
		client.SetBaseURL(d.APIBaseURL)
	}
	return client
}

func (d *Driver) createDevice() (*api.DeviceCreateResponse, error) {
	deviceCreateConfiguration := &api.DeviceCreateConfiguration{
		CPUCores:     d.CPUCores,
		DiskSize:     d.DiskSize,
		DisplayName:  d.MachineName,
		Hostname:     d.MachineName,
		KubernetesID: d.KubernetesID,
		Memory:       d.Memory,
		Password:     d.DevicePassword,
		SwapDiskSize: d.SwapDiskSize,
	}

	log.Debugf("Creating Xelon device with configuration: %+v", deviceCreateConfiguration)

	client := d.getClient()
	deviceCreateResponse, _, err := client.Devices.Create(deviceCreateConfiguration)
	if err != nil {
		return deviceCreateResponse, err
	}

	return deviceCreateResponse, nil
}

func (d *Driver) addSSHKey(localVMID string) error {
	d.SSHKeyPath = d.GetSSHKeyPath()

	if err := ssh.GenerateSSHKey(d.SSHKeyPath); err != nil {
		return err
	}

	publicKey, err := ioutil.ReadFile(d.SSHKeyPath + ".pub")
	if err != nil {
		return err
	}

	sshCreateConfiguration := &api.SSHAddRequest{
		Name:   d.MachineName,
		SSHKey: string(publicKey),
	}
	_, err = d.getClient().SSHs.Add(localVMID, sshCreateConfiguration)
	if err != nil {
		return err
	}

	return nil
}

func (d *Driver) startDevice() error {
	client := d.getClient()

	log.Debug("Checking device state...")
	deviceRoot, _, err := client.Devices.Get(d.TenantID, d.LocalVMID)
	if err != nil {
		return err
	}
	device := deviceRoot.Device

	if device.Powerstate == true && device.LocalVMDetails.State == 1 {
		log.Debug("Device is already running")
		return nil
	}

	log.Debug("Starting Xelon device...")
	_, err = client.Devices.Start(d.LocalVMID)
	if err != nil {
		return err
	}

	return nil
}

func (d *Driver) stopDevice() error {
	client := d.getClient()

	log.Debug("Checking device state...")
	deviceRoot, _, err := client.Devices.Get(d.TenantID, d.LocalVMID)
	if err != nil {
		return err
	}
	device := deviceRoot.Device

	if device.Powerstate == false {
		log.Debug("Device is already stopped")
		return nil
	}

	log.Debug("Stopping Xelon device...")
	_, err = client.Devices.Stop(d.LocalVMID)
	if err != nil {
		return err
	}

	log.Debug("Waiting until device is stopped...")
	for {
		deviceRoot, _, err := client.Devices.Get(d.TenantID, d.LocalVMID)
		if err != nil {
			return nil
		}
		if deviceRoot.Device.Powerstate == false {
			break
		}
		time.Sleep(1 * time.Second)
	}

	return nil
}

func randomDelay() {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(10)
	log.Debugf("random delay is %d", n)
	time.Sleep(time.Duration(n) * time.Second)
}
