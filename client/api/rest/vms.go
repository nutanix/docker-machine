package rest

import (
	"encoding/json"
	"fmt"

	"nutanix/client/api/mgmt"
)

func (c *NutanixRESTClient) GetVMInfo(vmId string) (*VMDTO, error) {
	respBytes, err := c.DoRequest("GET", fmt.Sprintf("%s%s", "/vms/", vmId), nil, nil)
	if err != nil {
		return nil, err
	}
	vmInfo := &VMDTO{}
	err = json.Unmarshal(respBytes, vmInfo)
	if err != nil {
		return nil, err
	}
	return vmInfo, nil
}

func (c *NutanixRESTClient) CreateVM(vmConfig *mgmt.VMCreateDTO) (*mgmt.ReturnValueDTO, error) {
	respBytes, err := c.DoRequest("POST", "/vms/", nil, vmConfig)
	if err != nil {
		return nil, err
	}
	taskDO := &mgmt.ReturnValueDTO{}
	err = json.Unmarshal(respBytes, taskDO)
	if err != nil {
		return nil, err
	}
	return taskDO, nil
}
