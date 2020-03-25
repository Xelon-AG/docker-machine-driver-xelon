package api

import (
	"context"
	"fmt"
	"net/http"
)

const tenantBasePath = "tenant"

// TenantService handles communication with the user related methods of the Xelon API.
type TenantService service

type Tenant struct {
	TenantIdentifier string `json:"tenant_identifier"`
}

// Get provides information about user especially tenant id.
func (s *TenantService) Get() (*Tenant, *http.Response, error) {
	path := fmt.Sprintf("%s", tenantBasePath)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	tenant := new(Tenant)
	resp, err := s.client.Do(context.Background(), req, tenant)
	if err != nil {
		return nil, resp, err
	}

	return tenant, resp, nil
}
