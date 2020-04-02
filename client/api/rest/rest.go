package rest

import (
	"nutanix/client"
)

const (
	restPath = "PrismGateway/services/rest/v1"
)

type NutanixRESTClient struct {
	*client.NutanixClient
}

func NewNutanixRESTClient(hostname, username, password string) *NutanixRESTClient {
	c, _ := client.NewNutanixClient(hostname, username, password, restPath)
	return &NutanixRESTClient{c}
}
