package mgmt

import (
	"encoding/json"
)

func (c *NutanixMGMTClient) GetNetworkList() (*NetworkList, error) {
	respBytes, err := c.DoRequest("GET", "/networks/", nil, nil)
	if err != nil {
		return nil, err
	}
	networkList := &NetworkList{}
	err = json.Unmarshal(respBytes, networkList)
	if err != nil {
		return nil, err
	}
	return networkList, nil
}
