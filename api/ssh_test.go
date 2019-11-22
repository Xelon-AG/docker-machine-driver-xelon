package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSSHsService_Add_emptyLocalVMID(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, err := client.SSHs.Add("", nil)

	assert.Error(t, err)
	assert.Equal(t, ErrEmptyArgument.Error(), err.Error())
}

func TestSSHsService_Add_emptySSHAddRequest(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, err := client.SSHs.Add("localVMID", nil)

	assert.Error(t, err)
	assert.Equal(t, ErrEmptyPayloadNotAllowed.Error(), err.Error())
}
