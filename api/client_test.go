package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setup() (client *Client, mux *http.ServeMux, serverURL string, teardown func()) {
	mux = http.NewServeMux()

	apiHandler := http.NewServeMux()
	apiHandler.Handle("/api/service/", http.StripPrefix("/api/service", mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		_, _ = fmt.Fprintln(os.Stderr, "FAIL: Client.BaseURL path prefix is not preserved in the request URL:")
		_, _ = fmt.Fprintln(os.Stderr)
		_, _ = fmt.Fprintln(os.Stderr, "\t"+req.URL.String())
		_, _ = fmt.Fprintln(os.Stderr)
		_, _ = fmt.Fprintln(os.Stderr, "\tDid you accidentally use an absolute endpoint URL rather then relative?")
		http.Error(w, "client.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})
	server := httptest.NewServer(apiHandler)
	client = NewClient("token")
	client.BaseURL, _ = url.Parse(server.URL + "/api/service/")

	return client, mux, server.URL, server.Close
}

func TestClient_NewClient(t *testing.T) {
	client := NewClient("token")

	assert.Equal(t, "https://vdc.xelon.ch/api/service/", client.BaseURL.String())
	assert.Equal(t, "docker-machine-driver-xelon", client.UserAgent)
}

func TestClient_Do_httpError(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "{\"error\":\"bad request\"}", http.StatusBadRequest)
	})
	req, _ := client.NewRequest("GET", ".", nil)

	resp, err := client.Do(context.Background(), req, nil)

	assert.Error(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}
