package rest

type VMDTO struct {
	VMName                string   `json:"vmName,omitempty"`
	ControllerVM          bool     `json:"controllerVm,omitempty"`
	VMId                  string   `json:"vmId,omitempty"`
	MemoryCapacityInBytes int      `json:"memoryCapacityInBytes,omitempty"`
	ContainerIds          []string `json:"containerIds,omitempty"`
	VirtualNicIds         []string `json:"virtualNicIds,omitempty"`
	PowerState            string   `json:"powerState,omitempty"`
	//onDiskDedup (dto.appliance.configuration.ContainerDTO$OnDiskDedup, optional),
	Stats             map[string]string `json:"stats,omitempty"`
	NumNetworkAdapter int               `json:"numNetworkAdapters,omitempty"`
	VDiskNames        []string          `json:"vdiskNames,omitempty"`
	//fingerPrintOnWrite (dto.appliance.configuration.ContainerDTO$FingerPrintOnWrite, optional),
	//healthSummary (get.dto.health.check.HealthSummaryDTO, optional),
	UsageStats           map[string]string `json:"usageStats,omitempty"`
	ClusterUUID          string            `json:"clusterUuid,omitempty"`
	CPUReservedInHz      int               `json:"cpuReservedInHz,omitempty"`
	HostName             string            `json:"hostName,omitempty"`
	ConsistencyGroupName string            `json:"consistencyGroupName,omitempty"`
	Displayable          bool              `json:"displayable,omitempty"`
	ProtectionDomainName string            `json:"protectionDomainName,omitempty"`
	HostId               string            `json:"hostId,omitempty"`
	NutanixVirtualDisks  []string          `json:"nutanixVirtualDisks,omitempty"`
	VDiskFilePaths       []string          `json:"vdiskFilePaths,omitempty"`
	HypervisorType       string            `json:"hypervisorType,omitempty"`
	NumVCPUs             int               `json:"numVCpus,omitempty"`
	//alertSummary (get.dto.alerts.AlertSummaryDTO, optional),
	NutanixVirtualDiskIds         []string `json:"nutanixVirtualDiskIds,omitempty"`
	RunningOnNDFS                 bool     `json:"runningOnNdfs,omitempty"`
	DiskCapacityInBytes           int      `json:"diskCapacityInBytes,omitempty"`
	GuestOperatingSystem          string   `json:"guestOperatingSystem,omitempty"`
	IpAddresses                   []string `json:"ipAddresses,omitempty"`
	MemortReservedCapacityInBytes int      `json:"memoryReservedCapacityInBytes,omitempty"`
	AcropolisVM                   bool     `json:"acropolisVm,omitempty"`
	NonNDFSDetails                string   `json:"nonNdfsDetails,omitempty"`
}

type NutanixGuestToolsDTO struct {
	VMId         string            `json:"vmId,omitempty"`
	Enabled      bool              `json:"enabled,omitempty"`
	Applications map[string]string `json:"applications,omitempty"`
}
