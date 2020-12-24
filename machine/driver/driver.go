package driver

import (
	"fmt"
	"io/ioutil"
	"net"
	"time"

	"nutanix/client/api/mgmt"
	"nutanix/client/api/rest"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/ssh"
	"github.com/docker/machine/libmachine/state"
	gouuid "github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/terraform-providers/terraform-provider-nutanix/client"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
)

const (
	defaultVMMem = 1024
	defaultVCPUs = 1
	defaultCores = 1
)

type NutanixDriver struct {
	*drivers.BaseDriver
	Username string
	Password string
	Endpoint string
	Cluster  string
	VMVCPUs  int
	VMCores  int
	VMMem    int
	SSHPass  string
	VLAN     string
	Image    string
	VMId     string
}

func NewDriver(hostname, storePath string) *NutanixDriver {
	return &NutanixDriver{
		BaseDriver: &drivers.BaseDriver{
			MachineName: hostname,
			StorePath:   storePath,
		},
	}
}

func (d *NutanixDriver) Create() error {
	name := d.GetMachineName()

	configCreds := client.Credentials{
		URL:         fmt.Sprintf("%s:%s", c.Endpoint, c.Port),
		Endpoint:    d.Endpoint,
		Username:    d.Username,
		Password:    d.Password,
		Port:        d.Port,
		Insecure:    d.Insecure,
		SessionAuth: d.SessionAuth,
		ProxyURL:    d.ProxyURL,
	}

	c := mgmt.NewNutanixMGMTClient(d.Endpoint, d.Username, d.Password)
	r := rest.NewNutanixRESTClient(d.Endpoint, d.Username, d.Password)
	v3Client, err := v3.NewV3Client(configCreds)
	if err != nil {
		return nil, err
	}

	uuid := gouuid.New().String()

	vmConfig := &mgmt.VMCreateDTO{
		MemoryMB:              d.VMMem,
		Name:                  name,
		NumVcpus:              d.VMVCPUs,
		NumCoresPerVcpu:       d.VMCores,
		UUID:                  uuid,
		VMDisks:               []*mgmt.VMDiskDTO{},
		VMNics:                []*mgmt.VMNicSpecDTO{},
		VMCustomizationConfig: &mgmt.VMCustomizationConfigDTO{},
	}
	networks, err := c.GetNetworkList()
	if err != nil {
		log.Errorf("Error getting networks: [%v]", err)
		return err
	}

	for _, net := range networks.Entities {
		if net.Name == d.VLAN {
			n := &mgmt.VMNicSpecDTO{
				NetworkUUID: net.UUID,
			}
			vmConfig.VMNics = append(vmConfig.VMNics, n)
			break
		}
	}

	images, err := c.GetImageList()
	if err != nil {
		log.Errorf("Error getting images: [%v]", err)
		return err
	}

	for _, img := range images.Entities {
		if img.Name == d.Image {
			d := &mgmt.VMDiskDTO{
				VMDiskClone: &mgmt.VMDiskSpecCloneDTO{
					VMDiskUUID: img.VMDiskID,
				},
			}
			vmConfig.VMDisks = append(vmConfig.VMDisks, d)
			break
		}
	}

	err = ssh.GenerateSSHKey(d.GetSSHKeyPath())
	if err != nil {
		log.Errorf("Error generating ssh key")
		return err
	}

	pubKey, err := ioutil.ReadFile(fmt.Sprintf("%s.pub", d.GetSSHKeyPath()))
	if err != nil {
		log.Errorf("Error reading public key")
		return err
	}
	vmConfig.VMCustomizationConfig.DataSourceType = "CONFIG_DRIVE_V2"
	vmConfig.VMCustomizationConfig.Userdata = "#cloud-config\r\nusers:\r\n - name: root\r\n   ssh_authorized_keys:\r\n    - " + string(pubKey)

	_, err = r.CreateVM(vmConfig)
	if err != nil {
		log.Errorf("Error creating vm: [%v]", err)
		return err
	}

	vmId := uuid
	for i := 0; i < 1200; i++ {
		vmDTO, err := r.GetVMInfo(uuid)
		if err != nil || len(vmDTO.NutanixVirtualDisks) < (2) {
			<-time.After(1 * time.Second)
			continue
		}
		break
	}
	d.VMId = vmId

	taskDO, err := c.PowerOn(vmId)
	if err != nil {
		log.Errorf("Error powering vm on: [%v]", err)
		return err
	}

	_, err = c.Wait(taskDO.TaskUUID)
	if err != nil {
		log.Errorf("Error waiting for power on task to complete: [%v]", err)
		return err
	}

	var vmDTO *rest.VMDTO
	ipAddr := ""

	doneChan := make(chan bool, 1)
	errChan := make(chan error, 1)

	go func(doneChan chan bool, errChan chan error) {
		for {
			select {
			case <-doneChan:
				// used to stop the goroutine if needed
				break
			default:
			}
			var err error
			vmDTO, err = r.GetVMInfo(vmId)
			if err != nil {
				log.Errorf("Error getting vm data from rest api: [%v]", err)
				errChan <- err
				break
			}
			if len(vmDTO.IpAddresses) > 0 {
				ipAddr = vmDTO.IpAddresses[0]
				doneChan <- true
				break
			}
			<-time.After(5 * time.Second)
		}
	}(doneChan, errChan)

	select {
	case <-doneChan:
	case err := <-errChan:
		return err
	case <-time.After(5 * time.Minute):
		doneChan <- false //end the go routine looking for ip address
		return fmt.Errorf("Too many retries to wait for IP address.")
	}

	d.IPAddress = ipAddr

	log.Infof("Created Nutanix Host %s, IP: %s", name, d.IPAddress)
	return nil
}

func (d *NutanixDriver) DriverName() string {
	return "nutanix"
}

func (d *NutanixDriver) GetCreateFlags() []mcnflag.Flag {
	return []mcnflag.Flag{
		mcnflag.StringFlag{
			EnvVar: "NUTANIX_USERNAME",
			Name:   "nutanix-username",
			Usage:  "Nutanix management username",
		},
		mcnflag.StringFlag{
			EnvVar: "NUTANIX_PASSWORD",
			Name:   "nutanix-password",
			Usage:  "Nutanix management password",
		},
		mcnflag.StringFlag{
			EnvVar: "NUTANIX_ENDPOINT",
			Name:   "nutanix-endpoint",
			Usage:  "Nutanix management endpoint ip address/FQDN",
		},
		mcnflag.IntFlag{
			EnvVar: "NUTANIX_VM_MEM",
			Name:   "nutanix-vm-mem",
			Usage:  "Memory in MB of the VM to be created",
			Value:  defaultVMMem,
		},
		mcnflag.IntFlag{
			EnvVar: "NUTANIX_VM_CPUS",
			Name:   "nutanix-vm-cpus",
			Usage:  "Number of VCPUs of the VM to be created",
			Value:  defaultVCPUs,
		},
		mcnflag.IntFlag{
			EnvVar: "NUTANIX_VM_CORES",
			Name:   "nutanix-vm-cores",
			Usage:  "Number of cores per VCPU of the VM to be created",
			Value:  defaultCores,
		},
		mcnflag.StringFlag{
			EnvVar: "NUTANIX_VM_NETWORK",
			Name:   "nutanix-vm-network",
			Usage:  "The name of the network to attach to the newly created VM",
		},
		mcnflag.StringFlag{
			EnvVar: "NUTANIX_VM_IMAGE",
			Name:   "nutanix-vm-image",
			Usage:  "The name of the VM disk to clone from, for the newly created VM",
		},
	}
}

func (d *NutanixDriver) GetSSHHostname() (string, error) {
	return d.GetIP()
}

func (d *NutanixDriver) GetURL() (string, error) {
	ip, err := d.GetIP()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("tcp://%s", net.JoinHostPort(ip, "2376")), nil
}

func (d *NutanixDriver) GetState() (state.State, error) {
	m := mgmt.NewNutanixMGMTClient(d.Endpoint, d.Username, d.Password)
	vmInfoDTO, err := m.GetVMInfo(d.VMId)
	if err != nil {
		return state.Error, err
	}
	switch vmInfoDTO.State {
	case "on":
		return state.Running, nil
	case "off":
		return state.Stopped, nil
	}
	return state.None, nil
}

func (d *NutanixDriver) Kill() error {
	m := mgmt.NewNutanixMGMTClient(d.Endpoint, d.Username, d.Password)
	taskDO, err := m.PowerOff(d.VMId)
	if err != nil {
		return err
	}
	_, err = m.Wait(taskDO.TaskUUID)
	return err
}

func (d *NutanixDriver) Remove() error {
	m := mgmt.NewNutanixMGMTClient(d.Endpoint, d.Username, d.Password)
	taskDO, err := m.DeleteVM(d.VMId)
	if err != nil {
		return err
	}
	_, err = m.Wait(taskDO.TaskUUID)
	return err
}

func (d *NutanixDriver) Restart() error {
	err := d.Stop()
	if err != nil {
		return err
	}
	return d.Start()
}

func (d *NutanixDriver) SetConfigFromFlags(opts drivers.DriverOptions) error {
	d.Username = opts.String("nutanix-username")
	if d.Username == "" {
		return fmt.Errorf("nutanix-username cannot be empty")
	}
	d.Password = opts.String("nutanix-password")
	if d.Password == "" {
		return fmt.Errorf("nutanix-password cannot be empty")
	}
	d.Endpoint = opts.String("nutanix-endpoint")
	if d.Endpoint == "" {
		return fmt.Errorf("nutanix-endpoint cannot be empty")
	}
	d.VMMem = opts.Int("nutanix-vm-mem")
	d.VMVCPUs = opts.Int("nutanix-vm-cpus")
	d.VMCores = opts.Int("nutanix-vm-cores")
	d.VLAN = opts.String("nutanix-vm-network")
	if d.VLAN == "" {
		return fmt.Errorf("nutanix-vm-network cannot be empty")
	}
	d.Image = opts.String("nutanix-vm-image")
	if d.Image == "" {
		return fmt.Errorf("nutanix-vm-image cannot be empty")
	}
	return nil
}

func (d *NutanixDriver) Start() error {
	m := mgmt.NewNutanixMGMTClient(d.Endpoint, d.Username, d.Password)
	taskDO, err := m.PowerOn(d.VMId)
	if err != nil {
		return err
	}
	_, err = m.Wait(taskDO.TaskUUID)
	return err
}

func (d *NutanixDriver) Stop() error {
	m := mgmt.NewNutanixMGMTClient(d.Endpoint, d.Username, d.Password)
	taskDO, err := m.PowerOff(d.VMId)
	if err != nil {
		return err
	}
	_, err = m.Wait(taskDO.TaskUUID)
	return err
}
