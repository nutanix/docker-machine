package mgmt

import (
	"nutanix/client"
)

const (
	mgmtPath = "api/nutanix/v0.8"
)

type NutanixMGMTClient struct {
	*client.NutanixClient
}

func NewNutanixMGMTClient(hostname, username, password string) *NutanixMGMTClient {
	c, _ := client.NewNutanixClient(hostname, username, password, mgmtPath)
	return &NutanixMGMTClient{c}
}
