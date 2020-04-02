package tests

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"

	"nutanix/client/api/mgmt"
)

func TestLifecycle(t *testing.T) {
	testMachineName := fmt.Sprintf("testMachine%d", rand.Int())

	sshUser := os.Getenv("NUTANIX_SSH_USER")
	sshPass := os.Getenv("NUTANIX_SSH_PASS")
	nutanixUser := os.Getenv("NUTANIX_USER")
	nutanixPass := os.Getenv("NUTANIX_PASS")
	endpoint := os.Getenv("NUTANIX_ENDPOINT")

	vNicList := os.Getenv("NUTANIX_VNICS")
	vNics := strings.Split(vNicList, ",")

	vDiskList := os.Getenv("NUTANIX_VDISKS")
	vDisks := strings.Split(vDiskList, ",")

	m := MachineConfig{
		MachineName: testMachineName,
		SSHUser:     sshUser,
		SSHPass:     sshPass,
		NutanixUser: nutanixUser,
		NutanixPass: nutanixPass,
		Endpoint:    endpoint,
		VNICs:       vNics,
		VDisks:      vDisks,
	}

	out, err := m.CreateMachine()
	if err != nil {
		t.Error("creation of machine failed", err, string(out))
		t.Fail()
	}

	defer m.DeleteMachine()

	c := mgmt.NewNutanixMGMTClient(endpoint, nutanixUser, nutanixPass)
	vmList, err := c.GetVMList()
	if err != nil {
		t.Error("Failed to get VM list from nutanix", err)
		t.Fail()
	}

	found := false

	var vmOfInterest *mgmt.VMInfoDTO

	for _, vm := range vmList.Entities {
		if vm.Config.Name == testMachineName {
			found = true
			vmOfInterest = vm
			break
		}
	}

	if found == false {
		t.Error("Docker machine create did not succeed, could not find VM by name - ", testMachineName)
		t.Fail()
	}

	if vmOfInterest.State != "on" {
		t.Error("Docker machine create did not switch on the machine. found state = ", vmOfInterest.State)
		t.Fail()
	}

	if vmOfInterest.Config.MemoryMB != 1024 {
		t.Error("Docker machine did not set the correct default value for VM memory, found [", vmOfInterest.Config.MemoryMB, " expected [1024 MB]")
		t.Fail()
	}

	if vmOfInterest.Config.NumCoresPerVcpu != 1 {
		t.Error("Docker machine did not set the correct default value for VM cores, found [", vmOfInterest.Config.NumCoresPerVcpu, " expected [1]")
		t.Fail()
	}

	if vmOfInterest.Config.NumVcpus != 1 {
		t.Error("Docker machine did not set the correct default value for VM CPUs, found [", vmOfInterest.Config.NumVcpus, " expected [1]")
		t.Fail()
	}

	foundNics := []string{}
	for _, nic := range vmOfInterest.Config.VmNics {
		foundNics = append(foundNics, nic.NetworkUUID)
	}

	sort.Strings(foundNics)
	sort.Strings(vNics)

	if !reflect.DeepEqual(foundNics, vNics) {
		t.Error("Docker machine did not set the correct values of Nics")
		t.Fail()
	}

}
