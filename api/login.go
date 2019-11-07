package api

import (
	"context"
	"fmt"
	"net/http"
)

const loginBasePath = "login"

type LoginService service

// User represents a Xelon user.
type User struct {
	APIToken         string `json:"api_token"`
	FirstName        string `json:"firstname"`
	ID               int    `json:"id"`
	Surname          string `json:"surname"`
	TenantIdentifier string `json:"tenantIdentifier"`
}

type userRoot struct {
	User User `json:"user"`
}

// LoginUser authenticates a user in application.
func (s *LoginService) LoginUser() (*User, *http.Response, error) {
	path := fmt.Sprintf("%v?email=%v&password=%v", loginBasePath, s.client.Username, s.client.Password)

	req, err := s.client.NewRequest(http.MethodPost, path, nil)
	if err != nil {
		return nil, nil, err
	}

	userRoot := new(userRoot)
	resp, err := s.client.Do(context.Background(), req, userRoot)

	return &userRoot.User, resp, err
}
