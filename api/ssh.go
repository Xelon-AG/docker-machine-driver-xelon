package api

import (
	"context"
	"fmt"
	"net/http"
)

const sshBasePath = "ssh"

// SSHsService handles communication with the ssh related methods of the Xelon API.
type SSHsService service

type SSHCreateConfiguration struct {
	Name      string
	PublicKey string
}

// Add attaches new SSH to device with specific localvmid.
func (s *SSHsService) Add(localVMID string, config *SSHCreateConfiguration) (*http.Response, error) {
	if localVMID == "" {
		return nil, ErrEmptyArgument
	}
	if config == nil {
		return nil, ErrEmptyPayloadNotAllowed
	}

	path := fmt.Sprintf("%v/%v/%v/add?name=%v&ssh_key=%v", deviceBasePath, localVMID, sshBasePath, config.Name, config.PublicKey)

	req, err := s.client.NewRequest(http.MethodPost, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(context.Background(), req, nil)
}
