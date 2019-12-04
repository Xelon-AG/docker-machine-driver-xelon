package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDevicesService_Get_emptyTenantID(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, _, err := client.Devices.Get("", "localVMID")

	assert.Error(t, err)
	assert.Equal(t, ErrEmptyArgument.Error(), err.Error())
}

func TestDevicesService_Get_emptyLocalVMID(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, _, err := client.Devices.Get("tenantID", "")

	assert.Error(t, err)
	assert.Equal(t, ErrEmptyArgument.Error(), err.Error())
}

func TestDevicesService_Create_emptyDeviceCreateConfiguration(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, _, err := client.Devices.Create(nil)

	assert.Error(t, err)
	assert.Equal(t, ErrEmptyPayloadNotAllowed.Error(), err.Error())
}

func TestDevicesService_Delete_emptyLocalVMID(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, err := client.Devices.Delete("")

	assert.Error(t, err)
	assert.Equal(t, ErrEmptyArgument.Error(), err.Error())
}

func TestDevicesService_Start_emptyLocalVMID(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, err := client.Devices.Start("")

	assert.Error(t, err)
	assert.Equal(t, ErrEmptyArgument.Error(), err.Error())
}

func TestDevicesService_Stop_emptyLocalVMID(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	_, err := client.Devices.Stop("")

	assert.Error(t, err)
	assert.Equal(t, ErrEmptyArgument.Error(), err.Error())
}
