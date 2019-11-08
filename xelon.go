package xelon

import (
	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/state"
)

type Driver struct {
	*drivers.BaseDriver
}

func NewDriver(hostName, storePath string) *Driver {
	return &Driver{
		BaseDriver: &drivers.BaseDriver{
			MachineName: hostName,
			StorePath:   storePath,
		},
	}
}

func (d *Driver) Create() error {
	panic("implement me")
}

func (d *Driver) DriverName() string {
	return "xelon"
}

func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	panic("implement me")
}

func (d *Driver) GetSSHHostname() (string, error) {
	return d.GetIP()
}

func (d *Driver) GetSSHKeyPath() string {
	if d.SSHKeyPath == "" {
		d.SSHKeyPath = d.ResolveStorePath("id_rsa")
	}
	return d.SSHKeyPath
}

func (d *Driver) GetURL() (string, error) {
	panic("implement me")
}

func (d *Driver) GetState() (state.State, error) {
	panic("implement me")
}

func (d *Driver) Kill() error {
	panic("implement me")
}

func (d *Driver) PreCreateCheck() error {
	panic("implement me")
}

func (d *Driver) Remove() error {
	panic("implement me")
}

func (d *Driver) Restart() error {
	panic("implement me")
}

func (d *Driver) SetConfigFromFlags(opts drivers.DriverOptions) error {
	panic("implement me")
}

func (d *Driver) Start() error {
	panic("implement me")
}

func (d *Driver) Stop() error {
	panic("implement me")
}
