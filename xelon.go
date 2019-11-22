package xelon

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/ssh"

	"bitbucket.org/xelonvdc/docker-machine-driver-xelon/api"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/state"
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
	defaultTemplateID     = 7
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
	Password       string
	SwapDiskSize   int
	TemplateID     int
	TenantID       string
	Username       string
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
	user, _, err := client.LoginService.Login()
	if err != nil {
		return err
	}
	log.Debugf("User tenant id: %v", user.TenantIdentifier)
	d.TenantID = user.TenantIdentifier

	log.Info("Creating Xelon device...")
	deviceCreateResponse, err := d.createDevice()
	if err != nil {
		return err
	}
	log.Debugf("DeviceCreateResponse: %+v", deviceCreateResponse)

	d.LocalVMID = deviceCreateResponse.Device.LocalVMID
	d.IPAddress = deviceCreateResponse.IPs[0]

	// TODO: workaround until response (array -> element) issue will be fixed
	log.Debug("Workaround (array -> element json parsing). Wait 60 seconds...")
	time.Sleep(60 * time.Second)

	log.Info("Waiting until Xelon device will be provisioned...")
	for {
		device, _, err := client.Devices.Get(user.TenantIdentifier, deviceCreateResponse.Device.LocalVMID)
		if err != nil {
			return err
		}
		if device.Powerstate == true && device.LocalVMDetails.State == 1 {
			break
		}
		time.Sleep(1 * time.Second)
	}

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
		mcnflag.StringFlag{
			EnvVar: "XELON_PASSWORD",
			Name:   "xelon-password",
			Usage:  "Xelon password",
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
		mcnflag.IntFlag{
			EnvVar: "XELON_TEMPLATE_ID",
			Name:   "xelon-template-id",
			Usage:  "ISO template ID for the device",
			Value:  defaultTemplateID,
		},
		mcnflag.StringFlag{
			EnvVar: "XELON_USERNAME",
			Name:   "xelon-username",
			Usage:  "Xelon user mail",
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
	device, _, err := d.getClient().Devices.Get(d.TenantID, d.LocalVMID)
	if err != nil {
		return state.Error, err
	}

	if device == nil {
		return state.None, nil
	}

	if device.Powerstate == false {
		return state.Stopped, nil
	} else {
		if device.LocalVMDetails.State == 1 {
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
	if d.DevicePassword != "" {
		// TODO: add device password validation
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
	d.Password = opts.String("xelon-password")
	d.SSHPort = opts.Int("xelon-ssh-port")
	d.SSHUser = opts.String("xelon-ssh-user")
	d.SwapDiskSize = opts.Int("xelon-swap-disk-size")
	d.TemplateID = opts.Int("xelon-template-id")
	d.Username = opts.String("xelon-username")

	if d.Username == "" {
		return fmt.Errorf("xelon driver requires the --xelon-username option")
	}
	if d.Password == "" {
		return fmt.Errorf("xelon driver requires the --xelon-password option")
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
	client := api.NewClient(d.Username, d.Password)
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
		TemplateID:   d.TemplateID,
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
	device, _, err := client.Devices.Get(d.TenantID, d.LocalVMID)
	if err != nil {
		return err
	}

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
	device, _, err := client.Devices.Get(d.TenantID, d.LocalVMID)
	if err != nil {
		return err
	}

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
		device, _, err := client.Devices.Get(d.TenantID, d.LocalVMID)
		if err != nil {
			return nil
		}
		if device.Powerstate == false {
			break
		}
		time.Sleep(1 * time.Second)
	}

	return nil
}
