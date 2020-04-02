package tests

import (
	"fmt"
	"os/exec"
)

const (
	driver = "nutanix"
)

type MachineConfig struct {
	VMMem       int
	VMCPUs      int
	VMCores     int
	MachineName string
	VNICs       []string
	VDisks      []string
	SSHUser     string
	SSHPass     string
	NutanixUser string
	NutanixPass string
	Endpoint    string
}

func (c MachineConfig) CreateMachine() ([]byte, error) {
	vmMem := c.VMMem
	if vmMem == 0 {
		vmMem = 1024
	}
	vmCores := c.VMCores
	if vmCores == 0 {
		vmCores = 1
	}
	vmCPUs := c.VMCPUs
	if vmCPUs == 0 {
		vmCPUs = 1
	}
	args := []string{
		"create",
		"-d",
		"nutanix",
		"--nutanix-username",
		c.NutanixUser,
		"--nutanix-password",
		c.NutanixPass,
		"--nutanix-endpoint",
		c.Endpoint,
		"--nutanix-vm-mem",
		fmt.Sprintf("%d", vmMem),
		"--nutanix-vm-cpus",
		fmt.Sprintf("%d", vmCPUs),
		"--nutanix-vm-cores",
		fmt.Sprintf("%d", vmCores),
		"--nutanix-vm-ssh-username",
		c.SSHUser,
		"--nutanix-vm-ssh-password",
		c.SSHPass,
	}

	for _, nic := range c.VNICs {
		args = append(args, "--nutanix-vm-network-uuid", nic)
	}

	for _, disk := range c.VDisks {
		args = append(args, "--nutanix-vm-disk-uuid", disk)
	}

	args = append(args, c.MachineName)

	cmd := exec.Command("docker-machine", args...)

	fmt.Println(cmd.Args)

	return cmd.CombinedOutput()
}

func (c MachineConfig) DeleteMachine() ([]byte, error) {
	args := []string{
		"rm",
		c.MachineName,
	}
	cmd := exec.Command("docker-machine", args...)
	return cmd.CombinedOutput()
}

func (c MachineConfig) StartMachine() ([]byte, error) {
	args := []string{
		"start",
		c.MachineName,
	}
	cmd := exec.Command("docker-machine", args...)
	return cmd.CombinedOutput()
}

func (c MachineConfig) StopMachine() ([]byte, error) {
	args := []string{
		"stop",
		c.MachineName,
	}
	cmd := exec.Command("docker-machine", args...)
	return cmd.CombinedOutput()
}
