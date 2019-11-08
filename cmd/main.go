package main

import (
	"bitbucket.org/xelonvdc/docker-machine-driver-xelon"
	"github.com/docker/machine/libmachine/drivers/plugin"
)

func main() {
	plugin.RegisterDriver(xelon.NewDriver("", ""))
}
