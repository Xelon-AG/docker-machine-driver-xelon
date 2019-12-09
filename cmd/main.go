package main

import (
	"github.com/Xelon-AG/docker-machine-driver-xelon"
	"github.com/docker/machine/libmachine/drivers/plugin"
)

func main() {
	plugin.RegisterDriver(xelon.NewDriver("", ""))
}
