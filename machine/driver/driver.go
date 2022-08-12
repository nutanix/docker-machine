package driver

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/ssh"
	"github.com/docker/machine/libmachine/state"
	"github.com/nutanix/docker-machine/utils"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	client "github.com/nutanix-cloud-native/prism-go-client/pkg/nutanix"
	v3 "github.com/nutanix-cloud-native/prism-go-client/pkg/nutanix/v3"
)

const (
	defaultVMMem = 2048
	defaultVCPUs = 2
	defaultCores = 1
)

// NutanixDriver driver structure
type NutanixDriver struct {
	*drivers.BaseDriver
	Endpoint         string
	Username         string
	Password         string
	Port             string
	Insecure         bool
	Cluster          string
	VMVCPUs          int
	VMCores          int
	VMCPUPassthrough bool
	VMMem            int
	SSHPass          string
	Subnet           []string
	Image            string
	ImageSize        int
	VMId             string
	SessionAuth      bool
	ProxyURL         string
	Categories       []string
	StorageContainer string
	DiskSize         int
	CloudInit        string
	SerialPort       bool
}

// NewDriver create new instance
func NewDriver(hostname, storePath string) *NutanixDriver {
	return &NutanixDriver{
		BaseDriver: &drivers.BaseDriver{
			MachineName: hostname,
			StorePath:   storePath,
		},
	}
}

// Create a host using the driver's config
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

	if d.VMCPUPassthrough {
		res.EnableCPUPassthrough = utils.BoolPtr(d.VMCPUPassthrough)
	}

	if d.SerialPort {
		SerialPort := &v3.VMSerialPort{
			Index:       utils.Int64Ptr(0),
			IsConnected: utils.BoolPtr(true),
		}

		res.SerialPortList = append(res.SerialPortList, SerialPort)
	}

	// Search target cluster
	c := &url.URL{Path: d.Cluster}
	encodedCluster := c.String()
	clusterFilter := fmt.Sprintf("name==%s", encodedCluster)
	clusters, err := conn.V3.ListAllCluster(clusterFilter)
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

	for index, subnet := range d.Subnet {
		// Trim extraneous whitespace
		d.Subnet[index] = strings.TrimSpace(subnet)
	}

	subnetFilter := ""

	for _, subnet := range d.Subnet {
		if len(subnetFilter) != 0 {
			subnetFilter += ","
		}

		t := &url.URL{Path: subnet}
		encodedSubnet := t.String()
		subnetFilter += fmt.Sprintf("name==%s", encodedSubnet)
	}

	subnets, err := conn.V3.ListAllSubnet(subnetFilter, getEmptyClientSideFilter())
	if err != nil {
		log.Errorf("Error getting subnets: [%v]", err)
		return err
	}

	for _, query := range d.Subnet {
		for _, subnet := range subnets.Entities {
			if *subnet.Status.Name == query && *subnet.Status.ClusterReference.UUID == *spec.ClusterReference.UUID {
				n := &v3.VMNic{
					SubnetReference: utils.BuildReference(*subnet.Metadata.UUID, "subnet"),
				}

				res.NicList = append(res.NicList, n)
				log.Infof("Subnet %s find with UUID: %s", *subnet.Status.Name, *subnet.Metadata.UUID)
				break
			}
		}
	}

	if len(res.NicList) < 1 {
		log.Errorf("Network %s not found in cluster %s", d.Subnet, d.Cluster)
		return fmt.Errorf("network %s not found in cluster %s", d.Subnet, d.Cluster)
	}

	if len(d.Categories) != 0 {
		log.Infof("Categories provided: %s", d.Categories)
		metadata.Categories = make(map[string]string)

		for _, group := range d.Categories {
			category := strings.Split(group, "=")

			if len(category) < 2 {
				log.Errorf("Malformed group %s", group)
				return fmt.Errorf("malformed group %s", group)
			}

			// Strip extraneous whitespace to make this more error tolerant
			category[0] = strings.TrimSpace(category[0])
			category[1] = strings.TrimSpace(category[1])

			metadata.Categories[category[0]] = category[1]
			log.Infof("Added category %s: %s", category[0], category[1])
		}
	}

	// Search image template
	i := &url.URL{Path: d.Image}
	encodedImage := i.String()
	imageFilter := fmt.Sprintf("name==%s", encodedImage)
	images, err := conn.V3.ListAllImage(imageFilter)
	if err != nil {
		log.Errorf("Error getting images: [%v]", err)
		return err
	}

	for _, image := range images.Entities {
		if *image.Status.Name == d.Image {

			log.Infof("Image %s find with UUID: %s", *image.Status.Name, *image.Metadata.UUID)

			if *image.Status.Resources.ImageType != "DISK_IMAGE" {
				log.Errorf("Image %s is not a disk template", d.Image)
				return fmt.Errorf("image %s is not a disk template", d.Image)
			}

			if d.ImageSize > 0 {
				newSize := int64(d.ImageSize * 1024)
				n := &v3.VMDisk{
					DataSourceReference: utils.BuildReference(*image.Metadata.UUID, "image"),
					DiskSizeMib:         &newSize,
				}
				res.DiskList = append(res.DiskList, n)
			} else {
				n := &v3.VMDisk{
					DataSourceReference: utils.BuildReference(*image.Metadata.UUID, "image"),
				}
				res.DiskList = append(res.DiskList, n)
			}

			break
		}
	}

	if len(res.DiskList) < 1 {
		log.Errorf("Image %s not found", d.Image)
		return fmt.Errorf("image %s not found", d.Image)
	}

	// Add additional disks
	if len(d.StorageContainer) != 0 && d.DiskSize > 0 {
		n := &v3.VMDisk{
			DiskSizeBytes: utils.Int64Ptr(int64(d.DiskSize) * 1024 * 1024 * 1024),
			StorageConfig: &v3.VMStorageConfig{
				StorageContainerReference: &v3.StorageContainerReference{
					Kind: "storage_container",
					UUID: d.StorageContainer,
				},
			},
		}

		res.DiskList = append(res.DiskList, n)
		log.Infof("Added disk with %d GiB on storage container with UUID: %s", d.DiskSize, d.StorageContainer)
	}

	// SSH Key generation
	err = ssh.GenerateSSHKey(d.GetSSHKeyPath())
	if err != nil {
		log.Errorf("Error generating ssh key")
		return err
	}

	pubKey, err := os.ReadFile(fmt.Sprintf("%s.pub", d.GetSSHKeyPath()))
	if err != nil {
		log.Errorf("Error reading public key")
		return err
	}

	log.Infof("SSH pub key ready (%s)", pubKey)

	// CloudInit preparation

	var userdata []byte

	if d.CloudInit != "" {
		t := yaml.Node{Kind: yaml.DocumentNode, HeadComment: "cloud-config"}

		if !strings.HasPrefix(d.CloudInit, "#cloud-config") {
			return errors.New("cloud-init syntax error")
		}

		err = yaml.Unmarshal([]byte(d.CloudInit), &t)
		if err != nil {
			log.Fatalf("Cloud-init syntax error: %v", err)
			return err
		}

		if t.Content == nil {
			log.Infof("Use default Cloud-init")
			userdata = []byte("#cloud-config\r\nusers:\r\n - name: root\r\n   ssh_authorized_keys:\r\n    - " + string(pubKey))
		} else {
			log.Infof("Cloud-init merge")

			usersNode := iterateNode(&t, "users")

			if usersNode == nil {
				rootNode := t.Content[0]
				rootNode.Content = append(rootNode.Content, buildScalarNodes("users")...)
				usersNode = &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
				rootNode.Content = append(rootNode.Content, usersNode)
			}

			rancherNode := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
			rancherNode.Content = append(rancherNode.Content, buildStringNodes("name", "root", "")...)
			rancherNode.Content = append(rancherNode.Content, buildStringNodes("sudo", "ALL=(ALL) NOPASSWD:ALL", "")...)
			rancherNode.Content = append(rancherNode.Content, buildScalarNodes("ssh-authorized-keys")...)

			sshSeqNode := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
			sshSeqNode.Content = append(sshSeqNode.Content, buildScalarNodes(string(pubKey))...)

			rancherNode.Content = append(rancherNode.Content, sshSeqNode)
			usersNode.Content = append(usersNode.Content, rancherNode)

			userdata, err = yaml.Marshal(&t)
			if err != nil {
				log.Fatal(err)
			}
			log.Infof(string(userdata))
		}
	} else {
		log.Infof("Use default Cloud-init")
		userdata = []byte("#cloud-config\r\nusers:\r\n - name: root\r\n   ssh_authorized_keys:\r\n    - " + string(pubKey))
	}

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

	log.Infof("waiting for vm %s (%s) to create: task %s", name, uuid, taskUUID)

	// Wait end of the task
waitTask:
	for i := 0; i < 60; i++ {
		resp, err := conn.V3.GetTask(taskUUID)
		if err != nil {
			log.Errorf("Error getting task: [%v]", err)
			return err
		}

		switch *resp.Status {
		case "SUCCEEDED":
			log.Infof("VM %s creation task succeeded", name)
			break waitTask
		case "FAILED":
			log.Errorf("Error creating vm: [%v]", *resp.ErrorDetail)
			return errors.New(*resp.ErrorDetail)
		}
		if i == 59 {
			log.Errorf("Timeout waiting for vm %s to create", name)
			return errors.New("timeout waiting for vm to create")
		}
		log.Infof("VM %s creation is in %s state", name, *resp.Status)
		<-time.After(5 * time.Second)

	}

	// Wait for the VM to be available
	for i := 0; i < 60; i++ {
		vmIntent, err := conn.V3.GetVM(uuid)
		if err != nil {
			log.Errorf("Error getting vm: [%v]", err)
			return err
		}

		nbDisks := len(spec.Resources.DiskList) + 1

		log.Infof("Waiting until at least %d disks are present...", nbDisks)

		if len(vmIntent.Spec.Resources.DiskList) == nbDisks {
			break
		}

		if i == 59 {
			log.Errorf("Timeout waiting for vm %s to be available", name)
			return errors.New("timeout waiting for vm to be available")
		} else {
			log.Infof("Waiting VM %s availability", name)
			<-time.After(5 * time.Second)
			continue
		}

	}
	d.VMId = uuid

	log.Infof("VM %s successfully created", name)

	// Wait for the VM obtain an IP address
	for i := 0; i < 60; i++ {
		vmInfo, err := conn.V3.GetVM(uuid)
		if err != nil {
			log.Errorf("Error getting vm: [%v]", err)
			return err
		}

		if len(vmInfo.Status.Resources.NicList[0].IPEndpointList) != 0 {
			d.IPAddress = *vmInfo.Status.Resources.NicList[0].IPEndpointList[0].IP
			break
		}

		if i == 59 {
			log.Errorf("Timeout waiting for vm %s to obtain an IP address", name)
			return errors.New("timeout waiting for vm to obtain an IP address")
		} else {
			log.Infof("Waiting VM %s ip configuration", name)
			<-time.After(5 * time.Second)
			continue
		}
	}

	log.Infof("VM %s configured with ip address %s", name, d.IPAddress)

	return nil
}

// DriverName returns the name of the driver
func (d *NutanixDriver) DriverName() string {
	return "nutanix"
}

// GetCreateFlags returns the mcnflag.Flag slice representing the flags
// that can be set, their descriptions and defaults.
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
		mcnflag.BoolFlag{
			EnvVar: "NUTANIX_VM_CPU_PASSTHROUGH",
			Name:   "nutanix-vm-cpu-passthrough",
			Usage:  "Enable passthrough the host's CPU features to the newly created VM",
		},
		mcnflag.StringSliceFlag{
			Name:  "nutanix-vm-network",
			Usage: "The name of the network to attach to the newly created VM",
		},
		mcnflag.StringFlag{
			EnvVar: "NUTANIX_VM_IMAGE",
			Name:   "nutanix-vm-image",
			Usage:  "The name of the VM disk to clone from, for the newly created VM",
		},
		mcnflag.IntFlag{
			EnvVar: "NUTANIX_VM_IMAGE_SIZE",
			Name:   "nutanix-vm-image-size",
			Usage:  "Increase the size of the template image",
			Value:  0,
		},
		mcnflag.StringSliceFlag{
			Name:  "nutanix-vm-categories",
			Usage: "The name of the categories who will be applied to the newly created VM",
		},
		mcnflag.StringFlag{
			EnvVar: "NUTANIX_STORAGE_CONTAINER",
			Name:   "nutanix-storage-container",
			Usage:  "The UUID of the storage container",
			Value:  "",
		},
		mcnflag.IntFlag{
			EnvVar: "NUTANIX_DISK_SIZE",
			Name:   "nutanix-disk-size",
			Usage:  "The size of the attached disk",
			Value:  0,
		},
		mcnflag.StringFlag{
			EnvVar: "NUTANIX_CLOUD_INIT",
			Name:   "nutanix-cloud-init",
			Usage:  "Cloud-init configuration",
		},
		mcnflag.BoolFlag{
			EnvVar: "NUTANIX_VM_SERIAL_PORT",
			Name:   "nutanix-vm-serial-port",
			Usage:  "Attach a serial port to the newly created VM (type Null)",
		},
	}
}

// GetSSHHostname returns hostname for use with ssh
func (d *NutanixDriver) GetSSHHostname() (string, error) {
	return d.GetIP()
}

// GetURL returns a Docker compatible host URL for connecting to this host
func (d *NutanixDriver) GetURL() (string, error) {
	ip, err := d.GetIP()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("tcp://%s", net.JoinHostPort(ip, "2376")), nil
}

// GetState returns the state that the host is in (running, stopped, etc)
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

// Kill stops a host forcefully
func (d *NutanixDriver) Kill() error {
	return d.Stop()
}

// Remove a host
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

// Restart a host. This may just call Stop(); Start() if the provider does not
// have any special restart behaviour.
func (d *NutanixDriver) Restart() error {
	err := d.Stop()
	if err != nil {
		return err
	}
	return d.Start()
}

// SetConfigFromFlags configures the driver with the object that was returned
// by RegisterCreateFlags
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

	d.Categories = opts.StringSlice("nutanix-vm-categories")

	d.Cluster = opts.String("nutanix-cluster")
	if d.Cluster == "" {
		return fmt.Errorf("nutanix-cluster cannot be empty")
	}

	d.DiskSize = opts.Int("nutanix-disk-size")
	d.StorageContainer = opts.String("nutanix-storage-container")

	d.VMMem = opts.Int("nutanix-vm-mem")
	d.VMVCPUs = opts.Int("nutanix-vm-cpus")
	d.VMCores = opts.Int("nutanix-vm-cores")

	d.VMCPUPassthrough = opts.Bool("nutanix-vm-cpu-passthrough")

	d.Subnet = opts.StringSlice("nutanix-vm-network")
	if len(d.Subnet) == 0 {
		return fmt.Errorf("nutanix-vm-network cannot be empty")
	}
	d.Image = opts.String("nutanix-vm-image")
	if d.Image == "" {
		return fmt.Errorf("nutanix-vm-image cannot be empty")
	}
	d.ImageSize = opts.Int("nutanix-vm-image-size")
	d.CloudInit = opts.String("nutanix-cloud-init")
	d.SerialPort = opts.Bool("nutanix-vm-serial-port")
	return nil
}

// Start a host
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

// Stop a host gracefully
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

func getEmptyClientSideFilter() []*client.AdditionalFilter {
	return make([]*client.AdditionalFilter, 0)
}
