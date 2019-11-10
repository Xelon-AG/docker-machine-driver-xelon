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
	CPU        int  `json:"cpu"`
	Powerstate bool `json:"powerstate"`
	RAM        int  `json:"ram"`
}

type deviceRoot struct {
	Device Device `json:"device"`
}

// Get provides detailed information for a device identified by tenant and localvmid.
func (s *DevicesService) Get(tenant, localVMID string) (*Device, *http.Response, error) {
	if tenant == "" || localVMID == "" {
		return nil, nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("device?tenant=%v&localvmid=%v", tenant, localVMID)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	deviceRoot := new(deviceRoot)
	resp, err := s.client.Do(context.Background(), req, deviceRoot)
	if err != nil {
		return nil, resp, err
	}

	return &deviceRoot.Device, resp, nil
}
