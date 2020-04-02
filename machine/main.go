package main

import (
	"github.com/docker/machine/libmachine/drivers/plugin"
	nutanix "nutanix/machine/driver"
)

func main() {
	plugin.RegisterDriver(nutanix.NewDriver("", ""))
}
