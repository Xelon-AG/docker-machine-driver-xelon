package xelon

import (
	"testing"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/stretchr/testify/assert"
)

func TestDriver_PreCreateCheck_MissingPassword(t *testing.T) {
	driver := NewDriver("default", "path")
	flags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			"xelon-password": "",
			"xelon-username": "user@xelon.ch",
		},
		CreateFlags: driver.GetCreateFlags(),
	}

	err := driver.SetConfigFromFlags(flags)
	assert.Error(t, err)
}

func TestDriver_PreCreateCheck_MissingUsername(t *testing.T) {
	driver := NewDriver("default", "path")
	flags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			"xelon-password": "password",
			"xelon-username": "",
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
			"xelon-password": "password",
			"xelon-username": "user@xelon.ch",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	_ = driver.SetConfigFromFlags(flags)

	err := driver.PreCreateCheck()
	assert.Error(t, err)
}
