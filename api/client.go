package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultBaseURL   = "https://vdc.xelon.ch/api/service/"
	defaultUserAgent = "docker-machine-driver-xelon"
)

// A Client manages communication with the Xelon API.
type Client struct {
	client *http.Client // HTTP client used to communicate with the API.

	BaseURL   *url.URL // Base URL for API requests. BaseURL should always be specified with a trailing slash.
	UserAgent string   // User agent used when communicating with Xelon API.
	Token     string   // Token for Xelon API.

	common service // Reuse a single struct instead of allocating one for each service on the heap.

	Devices *DevicesService
	SSHs    *SSHsService
	Tenant  *TenantService
}

type service struct {
	client *Client
}

// NewClient returns a new Xelon API client. To use API methods provide the token.
func NewClient(token string) *Client {
	httpClient := &http.Client{
		Timeout: time.Second * 15,
	}
	c := &Client{
		client:    httpClient,
		UserAgent: defaultUserAgent,
		Token:     token,
	}
	c.SetBaseURL(defaultBaseURL)
	c.common.client = c

	c.Devices = (*DevicesService)(&c.common)
	c.SSHs = (*SSHsService)(&c.common)
	c.Tenant = (*TenantService)(&c.common)

	return c
}

// SetBaseURL overrides the default BaseURL.
func (c *Client) SetBaseURL(baseURL string) {
	parsedURL, _ := url.Parse(baseURL)
	c.BaseURL = parsedURL
}

// NewRequest creates an API request. A relative URL can be provided in urlStr, in which case it is resolved
// relative to the BaseURL of the Client. Relative URLs should always be specified without a preceding slash.
// If specified, the value pointed to by body is JSON encoded and included as the request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a traling slash, but %q does not", c.BaseURL)
	}
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if body != nil {
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if req.Header.Get("Authorization") == "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	return req, nil
}

// Do sends an API request and returns the API response. The API response is JSON decoded and stored in
// the value pointed to by v, or returned as an error if an API error has occurred.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	req = req.WithContext(ctx)
	resp, err := c.client.Do(req)
	if err != nil {
		// if we got an error, and the context has been canceled, the context's error is more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// if the error type is *url.Error, sanitize its URL before returning.
		if e, ok := err.(*url.Error); ok {
			if uri, err := url.Parse(e.URL); err == nil {
				e.URL = sanitizeURL(uri).String()
				return nil, e
			}
		}

		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	err = CheckResponse(resp)
	if err != nil {
		return resp, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				return nil, err
			}
		} else {
			decodedErr := json.NewDecoder(resp.Body).Decode(v)
			if decodedErr == io.EOF {
				// ignore EOF errors caused by empty response body
				decodedErr = nil
			}
			if decodedErr != nil {
				err = decodedErr
			}
		}
	}

	return resp, err
}

// CheckResponse checks the API response for errors, and returns them if present. A response is considered
// an error if it has a status code outside the 200 range.
func CheckResponse(resp *http.Response) error {
	if code := resp.StatusCode; code >= 200 && code <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: resp}
	data, err := ioutil.ReadAll(resp.Body)
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, &errorResponse.ErrorElement)
		if err != nil {
			return err
		}
	}
	return errorResponse
}

// sanitizeURL redacts the password parameter from the URL which may be exposed by the user.
func sanitizeURL(uri *url.URL) *url.URL {
	if uri == nil {
		return nil
	}
	params := uri.Query()
	if len(params.Get("password")) > 0 {
		params.Set("password", "REDACTED")
		uri.RawQuery = params.Encode()
	}

	return uri
}
