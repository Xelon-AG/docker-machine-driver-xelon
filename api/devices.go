package api

import (
	"context"
	"fmt"
	"net/http"
)

const deviceBasePath = "vmlist"

// DevicesService handles communication with the devices related methods of the Xelon API.
type DevicesService service

// Device represents a Xelon device.
type Device struct {
	CPU            int            `json:"cpu"`
	LocalVMDetails LocalVMDetails `json:"localvmdetails,omitempty"`
	Networks       []Network      `json:"networks,omitempty"`
	Powerstate     bool           `json:"powerstate"`
	RAM            int            `json:"ram"`
}

type ToolsStatus struct {
	RunningStatus string `json:"runningStatus,omitempty"`
	Version       string `json:"version,omitempty"`
	ToolsStatus   bool   `json:"toolsStatus,omitempty"`
}

type Network struct {
	IPAddress  string `json:"ip,omitempty"`
	Label      string `json:"label,omitempty"`
	MacAddress string `json:"macAddress,omitempty"`
}

type DeviceCreateConfiguration struct {
	CPUCores     int
	DiskSize     int
	DisplayName  string
	Hostname     string
	KubernetesID string
	Memory       int
	Password     string
	SwapDiskSize int
}

type DeviceCreateResponse struct {
	Device LocalVMDetails `json:"device,omitempty"`
	IPs    []string       `json:"ips,omitempty"`
}

type LocalVMDetails struct {
	CreatedAt     string   `json:"created_at"`
	HVSystemID    int      `json:"hv_system_id"`
	ISOMounted    string   `json:"iso_mounted,omitempty"`
	LocalVMID     string   `json:"localvmid"`
	SSHKeys       []SSHKey `json:"ssh_keys,omitempty"`
	State         int      `json:"state"`
	TemplateID    int      `json:"template_id"`
	UpdatedAt     string   `json:"updated_at"`
	UserID        int      `json:"user_id"`
	VMDisplayName string   `json:"vmdisplayname"`
	VMHostname    string   `json:"vmhostname"`
}

type DeviceRoot struct {
	ToolsStatus ToolsStatus `json:"toolsStatus,omitempty"`
	Device      Device      `json:"device,omitempty"`
}

// Get provides detailed information for a device identified by tenant and localvmid.
func (s *DevicesService) Get(tenantID, localVMID string) (*DeviceRoot, *http.Response, error) {
	if tenantID == "" || localVMID == "" {
		return nil, nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("device?tenant=%v&localvmid=%v", tenantID, localVMID)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	deviceRoot := new(DeviceRoot)
	resp, err := s.client.Do(context.Background(), req, deviceRoot)
	if err != nil {
		return nil, resp, err
	}

	return deviceRoot, resp, nil
}

// Create makes a new device with given parameters.
func (s *DevicesService) Create(config *DeviceCreateConfiguration) (*DeviceCreateResponse, *http.Response, error) {
	if config == nil {
		return nil, nil, ErrEmptyPayloadNotAllowed
	}

	path := fmt.Sprintf("%v/create?cpucores=%v&disksize=%v&displayname=%v&hostname=%v&kubernetes_id=%v&memory=%v&password=%v&swapdisksize=%v",
		deviceBasePath, config.CPUCores, config.DiskSize, config.DisplayName, config.Hostname,
		config.KubernetesID, config.Memory, config.Password, config.SwapDiskSize)

	req, err := s.client.NewRequest(http.MethodPost, path, nil)
	if err != nil {
		return nil, nil, err
	}

	deviceCreateResponse := new(DeviceCreateResponse)
	resp, err := s.client.Do(context.Background(), req, deviceCreateResponse)
	if err != nil {
		return nil, resp, err
	}

	return deviceCreateResponse, resp, nil
}

// Delete removes a server.
func (s *DevicesService) Delete(localVMID string) (*http.Response, error) {
	if localVMID == "" {
		return nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/%v?password=%v", deviceBasePath, localVMID, s.client.Password)

	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(context.Background(), req, nil)
}

// Start starts a server with specific localvmid.
func (s *DevicesService) Start(localVMID string) (*http.Response, error) {
	if localVMID == "" {
		return nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/%v/startserver", deviceBasePath, localVMID)

	req, err := s.client.NewRequest(http.MethodPost, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(context.Background(), req, nil)
}

// Start stops a server with specific localvmid.
func (s *DevicesService) Stop(localVMID string) (*http.Response, error) {
	if localVMID == "" {
		return nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/%v/stopserver", deviceBasePath, localVMID)

	req, err := s.client.NewRequest(http.MethodPost, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(context.Background(), req, nil)
}
