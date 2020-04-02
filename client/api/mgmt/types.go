package mgmt

type ReturnValueDTO struct {
	TaskUUID string `json:"taskUuid"`
}

type TaskPollResultDTO struct {
	TimedOut       bool     `json:"timedOut,omitempty"`
	IsUnrecognized bool     `json:"isUnrecognized,omitempty"`
	TaskInfo       *TaskDTO `json:"taskInfo,omitempty"`
}

type TaskDTO struct {
	ParentTaskUUID     string           `json:"parentTaskUuid,omitempty"`
	MetaResponse       *MetaResponseDTO `json:"metaResponse,omitempty"`
	Message            string           `json:"message,omitempty"`
	UUID               string           `json:"uuid,omitempty"`
	CreateTime         int              `json:"createTime,omitempty"`
	PercentageComplete int              `json:"percentageComplete,omitempty"`
	EntityList         []*EntityIdDTO   `json:"entityList,omitempty"`
	CompleteTime       int              `json:"completeTime,omitempty"`
	//progressStatus TaskDTO$Status
	ProgressStatus  string `json:"progressStatus,omitempty"`
	StartTime       int    `json:"startTime,omitempty"`
	LastUpdatedTime int    `json:"lastUpdatedTime,omitempty"`
	//operationType TaskDTO$OperationType
	OperationType   string          `json:"operationType,omitempty"`
	SubtaskUUIDList []string        `json:"subtaskUuidList,omitempty"`
	MetaRequest     *MetaRequestDTO `json:"metaRequest,omitempty"`
}

type MetaResponseDTO struct {
	ErrorDetail string `json:"errorDetail,omitempty"`
	Error       string `json:"error,omitempty"`
}

type EntityIdDTO struct {
	//entityType EntityIdDTO$Entity
	EntityType string `json:"entityType,omitempty"`
	EntityName string `json:"entityName,omitempty"`
	UUID       string `json:"uuid,omitempty"`
}

type MetaRequestDTO struct {
	MethodName string `json:"methodName,omitempty"`
}

type VMDiskCreateDTO struct {
	Disks              []*VMDiskDTO `json:"disks,omitempty"`
	VMLogicalTimestamp int          `json:"vmLogicalTimestamp,omitempty"`
}

type VMDiskDTO struct {
	VMDiskCreate      *VMDiskSpecCreateDTO `json:"vmDiskCreate,omitempty"`
	DiskAddress       *VMDiskAddressDTO    `json:"diskAddress,omitempty"`
	IsScsiPassThrough bool                 `json:"isScsiPassThrough,omitempty"`
	VMDiskClone       *VMDiskSpecCloneDTO  `json:"vmDiskClone,omitempty"`
	IsEmpty           bool                 `json:"isEmpty,omitempty"`
	IsCdrom           bool                 `json:"isCdrom,omitempty"`
}

type VMDiskSpecCreateDTO struct {
	Size          int64  `json:"size,omitempty"`
	SizeMb        int64  `json:"sizeMb,omitempty"`
	ContainerName string `json:"containerName,omitempty"`
	ContainerId   int    `json:"containerId,omitempty"`
}

type VMDiskSpecCloneDTO struct {
	Ndfs_filepath string `json:"ndfs_filepath,omitempty"`
	MinimumSizeMb int64  `json:"minimumSizeMb,omitempty"`
	ImagePath     string `json:"imagePath,omitempty"`
	MinimumSize   int64  `json:"minimumSize,omitempty"`
	VMDiskUUID    string `json:"vmDiskUuid,omitempty"`
	VMdisk_uuid   string `json:"vmdisk_uuid,omitempty"`
}

type VMList struct {
	Metadata *Metadata    `json:"metadata,omitempty"`
	Entities []*VMInfoDTO `json:"entities,omitempty"`
}

type NetworkList struct {
	Metadata *Metadata           `json:"metadata,omitempty"`
	Entities []*NetworkConfigDTO `json:"entities,omitempty"`
}

type ProgressList struct {
	Metadata *Metadata            `json:"metadata,omitempty"`
	Entities []*ProgressStatusDTO `json:"entities,omitempty"`
}

type ProgressStatusDTO struct {
	PercentageCompleted int      `json:"percentageCompleted,omitempty"`
	EntityId            []string `json:"entityId,omitempty"`
}

type NetworkConfigDTO struct {
	Annotation       string `json:"annotation,omitempty"`
	LogicalTimestamp int    `json:"logicalTimestamp,omitempty"`
	Name             string `json:"name,omitempty"`
	UUID             string `json:"uuid,omitempty"`
	VLANID           int    `json:"vlanId,omitempty"`
	VSwitchName      string `json:"vswitchName,omitempty"`
}

type ImageList struct {
	Metadata *Metadata       `json:"metadata,omitempty"`
	Entities []*ImageInfoDTO `json:"entities,omitempty"`
}

type ImageInfoDTO struct {
	Name     string `json:"name,omitempty"`
	VMDiskID string `json:"vmDiskId,omitempty"`
}

type VMDiskList struct {
	Metadata *Metadata         `json:"metadata,omitempty"`
	Entities []VMDiskConfigDTO `json:"entities,omitempty"`
}

type Metadata struct {
	GrandTotalEntities int    `json:"grandTotalEntities,omitempty"`
	NextCursor         string `json:"nextCursor,omitempty"`
	SearchString       string `json:"searchString,omitempty"`
	StartIndex         int    `json:"startIndex,omitempty"`
	PreviousCursor     string `json:"previousCursor,omitempty"`
	FilterCriteria     string `json:"filterCriteria,omitempty"`
	EndIndex           int    `json:"endIndex,omitempty"`
	Count              int    `json:"count,omitempty"`
	Page               int    `json:"page,omitempty"`
	SortCriteria       string `json:"sortCriteria,omitempty"`
	TotalEntities      int    `json:"totalEntities,omitempty"`
}

type VMInfoDTO struct {
	LogicalTimestamp int    `json:"logicalTimestamp,omitempty"`
	HostUUID         string `json:"hostUuid,omitempty"`
	//State            VMInfoDTO.VMState `json:"state"`
	State  string       `json:"state,omitempty"`
	UUID   string       `json:"uuid,omitempty"`
	Config *VMConfigDTO `json:"config,omitempty"`
}

type VMCreateDTO struct {
	VMDisks               []*VMDiskDTO              `json:"vmDisks,omitempty"`
	MemoryMB              int                       `json:"memoryMb,omitempty"`
	VMNics                []*VMNicSpecDTO           `json:"vmNics,omitempty"`
	Name                  string                    `json:"name,omitempty"`
	Description           string                    `json:"description,omitempty"`
	HaPriority            int                       `json:"haPriority,omitempty"`
	Boot                  *BootConfigDTO            `json:"boot,omitempty"`
	NumCoresPerVcpu       int                       `json:"numCoresPerVcpu,omitempty"`
	NumVcpus              int                       `json:"numVcpus,omitempty"`
	VMCustomizationConfig *VMCustomizationConfigDTO `json:"vmCustomizationConfig,omitempty"`
	UUID                  string                    `json:"uuid,omitempty"`
}

type VMCustomizationConfigDTO struct {
	Userdata       string `json:"userdata,omitempty"`
	DataSourceType string `json:"datasourceType,omitempty"`
}

type VMConfigDTO struct {
	VmDisks         []*VMDiskConfigDTO `json:"vmDisks,omitempty"`
	MemoryMB        int                `json:"memoryMb,omitempty"`
	VmNics          []*VMNicSpecDTO    `json:"vmNics,omitempty"`
	Name            string             `json:"name,omitempty"`
	Description     string             `json:"description,omitempty"`
	HaPriority      int                `json:"haPriority,omitempty"`
	Boot            *BootConfigDTO     `json:"boot,omitempty"`
	NumCoresPerVcpu int                `json:"numCoresPerVcpu,omitempty"`
	NumVcpus        int                `json:"numVcpus,omitempty"`
}

type VMDiskConfigDTO struct {
	VmDiskSize        int64             `json:"vmDiskSize,omitempty"`
	IsEmpty           bool              `json:"isEmpty,omitempty"`
	VmDiskUUID        string            `json:"vmDiskUuid,omitempty"`
	Id                string            `json:"id,omitempty"`
	Addr              *VMDiskAddressDTO `json:"addr,omitempty"`
	ContainerID       int               `json:"containerId,omitempty"`
	IsCdrom           bool              `json:"isCdrom,omitempty"`
	IsSCSIPassthrough bool              `json:"isSCSIPassthrough,omitempty"`
}

type VMDiskAddressDTO struct {
	//DeviceBus   VMDiskAddressDTO.DeviceBus `json:"deviceBus"`
	DeviceBus   string `json:"deviceBus,omitempty"`
	DeviceIndex int    `json:"deviceIndex,omitempty"`
}

type VMNicSpecDTO struct {
	MacAddress         string `json:"macAddress,omitempty"`
	Model              string `json:"model,omitempty"`
	RequestedIpAddress string `json:"requestedIpAddress,omitempty"`
	NetworkUUID        string `json:"networkUuid,omitempty"`
}

type BootConfigDTO struct {
	//BootDeviceType BootConfigDTO.BootDeviceType `json:"bootDeviceType"`
	BootDeviceType string            `json:"bootDeviceType,omitempty"`
	DiskAddress    *VMDiskAddressDTO `json:"diskAddress,omitempty"`
	MacAddr        string            `json:"macAddr,omitempty"`
}

type VMPowerStateDTO struct {
	Transition string `json:"transition,omitempty"`
	HostUUID   string `json:"hostUuid,omitempty"`
	VMId       string `json:"vmid,omitempty"`
}

type ContainerList struct {
	Metadata *Metadata       `json:"metadata,omitempty"`
	Entities []*ContainerDTO `json:"entities,omitempty"`
}

type ContainerDTO struct {
	ID                            string            `json:"id,omitempty"`
	ReplicationFactor             int               `json:"replicationFactor,omitempty"`
	ErasureCode                   string            `json:"erasureCode,omitempty"`
	ErasureCodeDelaySecs          int               `json:"erasureCodeDelaySecs,omitempty"`
	CompressionEnabled            bool              `json:"compressionEnabled,omitempty"`
	OplogReplicationFactor        int               `json:"oplogReplicationFactor,omitempty"`
	CompressionDelayInSecs        int               `json:"compressionDelayInSecs,omitempty"`
	DownMigrateTimesInSecs        map[string]int    `json:"downMigrateTimesInSecs,omitempty"`
	RandomIoPreference            []string          `json:"randomIoPreference,omitempty"`
	OnDiskDedup                   string            `json:"onDiskDedup,omitempty"`
	Stats                         map[string]string `json:"stats,omitempty"`
	IlmPolicy                     string            `json:"ilmPolicy,omitempty"`
	AdvertisedCapacity            int               `json:"advertisedCapacity,omitempty"`
	FingerPrintOnWrite            string            `json:"fingerPrintOnWrite,omitempty"`
	TotalExplicitReservedCapacity int               `json:"totalExplicitReservedCapacity,omitempty"`
	StoragePoolId                 string            `json:"storagePoolId,omitempty"`
	NfsWhitelist                  []string          `json:"nfsWhitelist,omitempty"`
	UsageStats                    map[string]string `json:"usageStats,omitempty"`
	ClusterUuid                   string            `json:"clusterUuid,omitempty"`
	ContainerUuid                 string            `json:"containerUuid,omitempty"`
	SeqIoPreference               []string          `json:"seqIoPreference,omitempty"`
	MappedRemoteContainers        map[string]string `json:"mappedRemoteContainers,omitempty"`
	MarkedForRemoval              bool              `json:"markedForRemoval,omitempty"`
	MaxCapacity                   int               `json:"maxCapacity,omitempty"`
	VstoreNameList                []string          `json:"vstoreNameList,omitempty"`
	TotalImplicitReservedCapacity int               `json:"totalImplicitReservedCapacity,omitempty"`
	NfsWhitelistInherited         bool              `json:"nfsWhitelistInherited,omitempty"`
	Name                          string            `json:"name,omitempty"`

	//AlertSummary                  *AlertSummaryDTO  `json:"alertSummary,omitempty"`
	//HealthSummary                 *HealthSummaryDTO `json:"healthSummary,omitempty"`
}
