package mgmt

import (
	"encoding/json"
	"fmt"
)

func (c *NutanixMGMTClient) CreateVM(vmConfig *VMCreateDTO) (*ReturnValueDTO, error) {
	respBytes, err := c.DoRequest("POST", "/vms/", nil, vmConfig)
	if err != nil {
		return nil, err
	}
	taskDO := &ReturnValueDTO{}
	err = json.Unmarshal(respBytes, taskDO)
	if err != nil {
		return nil, err
	}
	return taskDO, nil
}

func (c *NutanixMGMTClient) GetVMList() (*VMList, error) {
	respBytes, err := c.DoRequest("GET", "/vms/", nil, nil)
	if err != nil {
		return nil, err
	}
	vmList := &VMList{}
	err = json.Unmarshal(respBytes, vmList)
	if err != nil {
		return nil, err
	}
	return vmList, nil
}

func (c *NutanixMGMTClient) GetVMInfo(vmId string) (*VMInfoDTO, error) {
	respBytes, err := c.DoRequest("GET", fmt.Sprintf("%s%s", "/vms/", vmId), nil, nil)
	if err != nil {
		return nil, err
	}
	vmInfo := &VMInfoDTO{}
	err = json.Unmarshal(respBytes, vmInfo)
	if err != nil {
		return nil, err
	}
	return vmInfo, nil
}

func (c *NutanixMGMTClient) PowerOn(vmId string) (*ReturnValueDTO, error) {
	powerOnState := &VMPowerStateDTO{
		Transition: "on",
	}
	respBytes, err := c.DoRequest("POST", fmt.Sprintf("%s%s/set_power_state/", "/vms/", vmId), nil, powerOnState)
	if err != nil {
		return nil, err
	}
	taskDO := &ReturnValueDTO{}
	err = json.Unmarshal(respBytes, taskDO)
	if err != nil {
		return nil, err
	}
	return taskDO, nil
}

func (c *NutanixMGMTClient) PowerOff(vmId string) (*ReturnValueDTO, error) {
	powerOffState := &VMPowerStateDTO{
		Transition: "off",
	}
	respBytes, err := c.DoRequest("POST", fmt.Sprintf("%s%s/set_power_state/", "/vms/", vmId), nil, powerOffState)
	if err != nil {
		return nil, err
	}
	taskDO := &ReturnValueDTO{}
	err = json.Unmarshal(respBytes, taskDO)
	if err != nil {
		return nil, err
	}
	return taskDO, nil
}

func (c *NutanixMGMTClient) DeleteVM(vmId string) (*ReturnValueDTO, error) {
	respBytes, err := c.DoRequest("DELETE", fmt.Sprintf("%s%s", "/vms/", vmId), nil, nil)
	if err != nil {
		return nil, err
	}
	taskDO := &ReturnValueDTO{}
	err = json.Unmarshal(respBytes, taskDO)
	if err != nil {
		return nil, err
	}
	return taskDO, nil
}
