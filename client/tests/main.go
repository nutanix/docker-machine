package main

import (
	"encoding/json"
	"fmt"

	"nutanix/client/api/mgmt"
	"nutanix/client/api/rest"
)

func main() {
	c := mgmt.NewNutanixMGMTClient("192.168.1.199:9440", "admin", "Suite219")

	vmConfig := &mgmt.VMCreateDTO{
		MemoryMB:        2048,
		Name:            "TestAPI",
		NumVcpus:        2,
		NumCoresPerVcpu: 2,
		VMDisks: []*mgmt.VMDiskDTO{
			{
				VMDiskClone: &mgmt.VMDiskSpecCloneDTO{
					VMDiskUUID: "a94e256c-10b2-424e-874c-5a6bc10fffe6",
				},
			},
		},
		VMNics: []*mgmt.VMNicSpecDTO{
			{
				NetworkUUID: "e06799b3-e909-4663-af6e-4d427fca8e64",
			},
		},
	}

	taskDO, err := c.CreateVM(vmConfig)
	if err != nil {
		fmt.Printf("Error creating vm: [%v]", err)
		return
	}

	taskResult, err := c.Wait(taskDO.TaskUUID)
	if err != nil {
		fmt.Printf("Error waiting for create task to complete: [%v]", err)
		return
	}

	vmId := taskResult.TaskInfo.EntityList[0].UUID

	taskDO, err = c.PowerOn(vmId)
	if err != nil {
		fmt.Printf("Error powering vm on: [%v]", err)
		return
	}

	taskResult, err = c.Wait(taskDO.TaskUUID)
	if err != nil {
		fmt.Printf("Error waiting for power on task to complete: [%v]", err)
		return
	}

	r := rest.NewNutanixRESTClient("192.168.1.199:9440", "admin", "Suite219")

	var vmDTO *rest.VMDTO
	ipAddr := ""

	for i := 0; i < 2000; i++ {
		vmDTO, err = r.GetVMInfo(vmId)
		if err != nil {
			fmt.Printf("Error getting vm data from rest api: [%v]", err)
			return
		}
		if len(vmDTO.IpAddresses) > 0 {
			ipAddr = vmDTO.IpAddresses[0]
			break
		}
	}
	ipAddr = ipAddr
	bytes, err := json.MarshalIndent(vmDTO, "  ", "    ")
	if err != nil {
		fmt.Printf("Error marshalling vm info from rest api: [%v]", err)
		return
	}

	fmt.Printf("%s\n", bytes)

}
