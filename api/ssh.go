package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
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

	escapedName := url.PathEscape(config.Name)
	escapedPublicKey := url.PathEscape(config.PublicKey)
	path := fmt.Sprintf("%v/%v/%v/add?name=%v&ssh_key=%v", deviceBasePath, localVMID, sshBasePath, escapedName, escapedPublicKey)

	req, err := s.client.NewRequest(http.MethodPost, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(context.Background(), req, nil)
}
