package xelon

import (
	"testing"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/stretchr/testify/assert"
)

func TestDriver_PreCreateCheck_MissingToken(t *testing.T) {
	driver := NewDriver("default", "path")
	flags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			"xelon-token": "",
		},
		CreateFlags: driver.GetCreateFlags(),
	}

	err := driver.SetConfigFromFlags(flags)
	assert.Error(t, err)
}

func TestDriver_PreCreateCheck_DevicePasswordLessThen6Characters(t *testing.T) {
	driver := NewDriver("default", "path")
	flags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			"xelon-device-password": "12345",
			"xelon-token": "token",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	_ = driver.SetConfigFromFlags(flags)

	err := driver.PreCreateCheck()
	assert.Error(t, err)
}
