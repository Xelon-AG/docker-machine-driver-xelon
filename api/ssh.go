package api

import (
	"context"
	"fmt"
	"net/http"
)

const sshBasePath = "ssh"

// SSHsService handles communication with the ssh related methods of the Xelon API.
type SSHsService service

type SSHKey struct {
	CreatedAt string `json:"created_at,omitempty"`
	DeleteAt  string `json:"deleted_at,omitempty"`
	ID        int    `json:"id"`
	Name      string `json:"name"`
	PublicKey string `json:"ssh_key"`
	UpdatedAt string `json:"updated_at,omitempty"`
	UserID    int    `json:"user_id,omitempty"`
	VMID      int    `json:"vm_id,omitempty"`
}

type SSHAddRequest struct {
	Name   string `json:"name"`
	SSHKey string `json:"ssh_key"`
}

// Add attaches new SSH to device with specific localvmid.
func (s *SSHsService) Add(localVMID string, sshAddRequest *SSHAddRequest) (*http.Response, error) {
	if localVMID == "" {
		return nil, ErrEmptyArgument
	}
	if sshAddRequest == nil {
		return nil, ErrEmptyPayloadNotAllowed
	}

	path := fmt.Sprintf("%v/%v/%v/add", deviceBasePath, localVMID, sshBasePath)

	req, err := s.client.NewRequest(http.MethodPost, path, sshAddRequest)
	if err != nil {
		return nil, err
	}

	return s.client.Do(context.Background(), req, nil)
}
