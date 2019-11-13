package api

import (
	"context"
	"fmt"
	"net/http"
)

const loginBasePath = "login"

// LoginService handles communication with the login related methods of the Xelon API.
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

// Login authenticates user into application.
func (s *LoginService) Login() (*User, *http.Response, error) {
	path := fmt.Sprintf("%v?email=%v&password=%v", loginBasePath, s.client.Username, s.client.Password)

	req, err := s.client.NewRequest(http.MethodPost, path, nil)
	if err != nil {
		return nil, nil, err
	}

	userRoot := new(userRoot)
	resp, err := s.client.Do(context.Background(), req, userRoot)
	if err != nil {
		return nil, resp, err
	}

	return &userRoot.User, resp, nil
}
