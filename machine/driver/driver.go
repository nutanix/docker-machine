package driver

import (
	"fmt"
	"io/ioutil"
	"net"
	"time"
	"encoding/base64"

	"nutanix/utils"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/ssh"
	"github.com/docker/machine/libmachine/state"
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
	Endpoint    string
	Username    string
	Password    string
	Port        string
	Insecure    bool
	Cluster     string
	VMVCPUs     int
	VMCores     int
	VMMem       int
	SSHPass     string
	Subnet      string
	Image       string
	VMId        string
	SessionAuth bool
	ProxyURL    string
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
		URL:         fmt.Sprintf("%s:%s", d.Endpoint, d.Port),
		Endpoint:    d.Endpoint,
		Username:    d.Username,
		Password:    d.Password,
		Port:        d.Port,
		Insecure:    d.Insecure,
		SessionAuth: d.SessionAuth,
		ProxyURL:    d.ProxyURL,
	}

	log.Infof("Connecting on: %s", configCreds.URL)

	conn, err := v3.NewV3Client(configCreds)
	if err != nil {
		return err
	}

	// Prepare VM creation request
	request := &v3.VMIntentInput{}
	spec := &v3.VM{}
	metadata := &v3.Metadata{}
	res := &v3.VMResources{}


	res.MemorySizeMib = utils.Int64Ptr(int64(d.VMMem))
	res.NumSockets = utils.Int64Ptr(int64(d.VMVCPUs))
	res.NumVcpusPerSocket = utils.Int64Ptr(int64(d.VMCores))

	// Search target cluster 
	clusters, err := conn.V3.ListAllCluster("")
	if err != nil {
		log.Errorf("Error getting clusters: [%v]", err)
		return err
	}

	for _, cluster := range clusters.Entities {
		if *cluster.Status.Name == d.Cluster {
			
			log.Infof("Cluster %s find with UUID: %s", *cluster.Status.Name, *cluster.Metadata.UUID)
			spec.ClusterReference = utils.BuildReference(*cluster.Metadata.UUID, "cluster")
			break
		}
	}

	// Search target subnet
	subnets, err := conn.V3.ListAllSubnet("")
	if err != nil {
		log.Errorf("Error getting subnets: [%v]", err)
		return err
	}

	for _, subnet := range subnets.Entities {
		if *subnet.Status.Name == d.Subnet {
			
			n := &v3.VMNic{
				SubnetReference: utils.BuildReference(*subnet.Metadata.UUID, "subnet"),
			}

			res.NicList = append(res.NicList, n)
			log.Infof("Subnet %s find with UUID: %s", *subnet.Status.Name, *subnet.Metadata.UUID)
			break
		}
	}

	if len(res.NicList) < 1 {
		log.Errorf("Network %s not found", d.Subnet)
		return fmt.Errorf("Network %s not found", d.Subnet)
	}


	// Search image template
	images, err := conn.V3.ListAllImage("")
	if err != nil {
		log.Errorf("Error getting images: [%v]", err)
		return err
	}

	for _, image := range images.Entities {
		if *image.Status.Name == d.Image {
			
			n := &v3.VMDisk{
				DataSourceReference: utils.BuildReference(*image.Metadata.UUID, "image"),
			}

			res.DiskList = append(res.DiskList, n)
			log.Infof("Image %s find with UUID: %s", *image.Status.Name, *image.Metadata.UUID)
			break
		}
	}

	if len(res.DiskList) < 1 {
		log.Errorf("Image %s not found", d.Image)
		return fmt.Errorf("Image %s not found", d.Image)
	}

	// SSH Key generation
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

	log.Infof("SSH pub key ready (%s)", pubKey)

	// CloudInit preparation

	// vmConfig.VMCustomizationConfig.DataSourceType = "CONFIG_DRIVE_V2"
	userdata := []byte("#cloud-config\r\nusers:\r\n - name: root\r\n   ssh_authorized_keys:\r\n    - " + string(pubKey))

	cloudInit := &v3.GuestCustomizationCloudInit{
		UserData: utils.StringPtr(base64.StdEncoding.EncodeToString(userdata)),
	}

	guestCustomization := &v3.GuestCustomization{
		CloudInit: cloudInit,
	}

	res.GuestCustomization = guestCustomization

	metadata.Kind = utils.StringPtr("vm")
	spec.Name = utils.StringPtr(name)
	spec.Description = utils.StringPtr("VM created by docker-image")
	res.PowerState = utils.StringPtr("ON")
	spec.Resources = res
	request.Metadata = metadata
	request.Spec = spec

	log.Infof("Launch VM creation")
	resp, err := conn.V3.CreateVM(request)
	if err != nil {
		log.Errorf("Error creating vm: [%v]", err)
		return err
	}

	uuid := *resp.Metadata.UUID
	taskUUID := resp.Status.ExecutionContext.TaskUUID.(string)

	log.Infof("waiting for vm (%s) to create: %s", uuid, taskUUID)

	// Wait for the VM to be available
	for i := 0; i < 1200; i++ {
		vmIntent, err := conn.V3.GetVM(uuid)
		if err != nil || len(vmIntent.Spec.Resources.DiskList) < (2) {
			<-time.After(1 * time.Second)
			continue
		}
		break
	}
	d.VMId = uuid

	log.Infof("VM %s succesfully created", name )

	var vmInfo *v3.VMIntentResponse
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
			vmInfo, err = conn.V3.GetVM(uuid)
			if err != nil {
				log.Errorf("Error getting vm data from rest api: [%v]", err)
				errChan <- err
				break
			}
			if len(vmInfo.Status.Resources.NicList[0].IPEndpointList) > 0 {
				ipAddr = *vmInfo.Status.Resources.NicList[0].IPEndpointList[0].IP
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
		mcnflag.StringFlag{
			EnvVar: "NUTANIX_PORT",
			Name:   "nutanix-port",
			Usage:  "Nutanix management endpoint port (default: 9440)",
			Value:  "9440",
		},
		mcnflag.BoolFlag{
			EnvVar: "NUTANIX_INSECURE",
			Name:   "nutanix-insecure",
			Usage:  "Explicitly allow the provider to perform \"insecure\" SSL requests",
		},
		mcnflag.StringFlag{
			EnvVar: "NUTANIX_CLUSTER",
			Name:   "nutanix-cluster",
			Usage:  "Nutanix Cluster to install VM on",
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

	configCreds := client.Credentials{
		URL:         fmt.Sprintf("%s:%s", d.Endpoint, d.Port),
		Endpoint:    d.Endpoint,
		Username:    d.Username,
		Password:    d.Password,
		Port:        d.Port,
		Insecure:    d.Insecure,
		SessionAuth: d.SessionAuth,
		ProxyURL:    d.ProxyURL,
	}

	log.Infof("Connecting on: %s", configCreds.URL)

	conn, err := v3.NewV3Client(configCreds)
	if err != nil {
		return state.Error, err
	}
	
	resp, err := conn.V3.GetVM(d.VMId)
	if err != nil {
		return state.Error, err
	}
	switch *resp.Status.Resources.PowerState {
	case "ON":
		return state.Running, nil
	case "OFF":
		return state.Stopped, nil
	}
	return state.None, nil
}

func (d *NutanixDriver) Kill() error {
	return d.Stop()
}

func (d *NutanixDriver) Remove() error {
	name := d.GetMachineName()

	configCreds := client.Credentials{
		URL:         fmt.Sprintf("%s:%s", d.Endpoint, d.Port),
		Endpoint:    d.Endpoint,
		Username:    d.Username,
		Password:    d.Password,
		Port:        d.Port,
		Insecure:    d.Insecure,
		SessionAuth: d.SessionAuth,
		ProxyURL:    d.ProxyURL,
	}

	log.Infof("Connecting on: %s", configCreds.URL)

	conn, err := v3.NewV3Client(configCreds)
	if err != nil {
		return err
	}
	resp, err := conn.V3.DeleteVM(d.VMId)
	if err != nil {
		return err
	}

	taskUUID := resp.Status.ExecutionContext.TaskUUID.(string)

	// Wait for the VM to be deleted
	for i := 0; i < 1200; i++ {
		resp, err := conn.V3.GetTask(taskUUID)
		if err != nil || *resp.Status != "SUCCEEDED" {
			<-time.After(1 * time.Second)
			continue
		}
		return err
	}
	return fmt.Errorf("unable to delete VM %s", name)
	
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
	d.Port = opts.String("nutanix-port")

	d.Insecure = opts.Bool("nutanix-insecure")

	d.Cluster = opts.String("nutanix-cluster")
	if d.Cluster == "" {
		return fmt.Errorf("nutanix-cluster cannot be empty")
	}

	d.VMMem = opts.Int("nutanix-vm-mem")
	d.VMVCPUs = opts.Int("nutanix-vm-cpus")
	d.VMCores = opts.Int("nutanix-vm-cores")
	d.Subnet = opts.String("nutanix-vm-network")
	if d.Subnet == "" {
		return fmt.Errorf("nutanix-vm-network cannot be empty")
	}
	d.Image = opts.String("nutanix-vm-image")
	if d.Image == "" {
		return fmt.Errorf("nutanix-vm-image cannot be empty")
	}
	return nil
}

func (d *NutanixDriver) Start() error {
	name := d.GetMachineName()

	configCreds := client.Credentials{
		URL:         fmt.Sprintf("%s:%s", d.Endpoint, d.Port),
		Endpoint:    d.Endpoint,
		Username:    d.Username,
		Password:    d.Password,
		Port:        d.Port,
		Insecure:    d.Insecure,
		SessionAuth: d.SessionAuth,
		ProxyURL:    d.ProxyURL,
	}

	log.Infof("Connecting on: %s", configCreds.URL)

	conn, err := v3.NewV3Client(configCreds)
	if err != nil {
		return err
	}

	vmResp, err := conn.V3.GetVM(d.VMId)
	if err != nil {
		return err
	}

	// Prepare VM update request
	request := &v3.VMIntentInput{}
	request.Spec = vmResp.Spec
	request.Metadata = vmResp.Metadata
	request.Spec.Resources.PowerState = utils.StringPtr("ON")
	
	resp, err := conn.V3.UpdateVM(d.VMId, request)
	if err != nil {
		return err
	}

	taskUUID := resp.Status.ExecutionContext.TaskUUID.(string)

	// Wait for the VM to be deleted
	for i := 0; i < 1200; i++ {
		resp, err := conn.V3.GetTask(taskUUID)
		if err != nil || *resp.Status != "SUCCEEDED" {
			<-time.After(1 * time.Second)
			continue
		}
		return err
	}
	return fmt.Errorf("unable to Start VM %s", name)
}

func (d *NutanixDriver) Stop() error {
	name := d.GetMachineName()

	configCreds := client.Credentials{
		URL:         fmt.Sprintf("%s:%s", d.Endpoint, d.Port),
		Endpoint:    d.Endpoint,
		Username:    d.Username,
		Password:    d.Password,
		Port:        d.Port,
		Insecure:    d.Insecure,
		SessionAuth: d.SessionAuth,
		ProxyURL:    d.ProxyURL,
	}

	log.Infof("Connecting on: %s", configCreds.URL)

	conn, err := v3.NewV3Client(configCreds)
	if err != nil {
		return err
	}

	vmResp, err := conn.V3.GetVM(d.VMId)
	if err != nil {
		return err
	}

	// Prepare VM update request
	request := &v3.VMIntentInput{}
	request.Spec = vmResp.Spec
	request.Metadata = vmResp.Metadata
	request.Spec.Resources.PowerState = utils.StringPtr("OFF")
	
	resp, err := conn.V3.UpdateVM(d.VMId, request)
	if err != nil {
		return err
	}

	taskUUID := resp.Status.ExecutionContext.TaskUUID.(string)

	// Wait for the VM to be deleted
	for i := 0; i < 1200; i++ {
		resp, err := conn.V3.GetTask(taskUUID)
		if err != nil || *resp.Status != "SUCCEEDED" {
			<-time.After(1 * time.Second)
			continue
		}
		return err
	}
	return fmt.Errorf("unable to Stop VM %s", name)
}

