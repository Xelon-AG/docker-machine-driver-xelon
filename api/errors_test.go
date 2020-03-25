package api

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrors_Error_messageFormat(t *testing.T) {
	errorResponse := &ErrorResponse{
		Response: &http.Response{
			Request: &http.Request{
				Method: http.MethodPost,
				URL:    &url.URL{Scheme: "https", Path: "vdc.xelon.ch/api/service"},
			},
			StatusCode: 401,
		},
		ErrorElement: ErrorElement{
			Error: "Unauthenticated user",
			Code:  401,
		},
	}
	expectedMessage := "POST https://vdc.xelon.ch/api/service: 401 {Error:Unauthenticated user Code:401}"

	assert.Error(t, errorResponse)
	assert.Equal(t, expectedMessage, errorResponse.Error())
}
