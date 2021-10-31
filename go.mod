module nutanix

go 1.14

replace github.com/docker/docker => github.com/docker/engine v17.12.0-ce-rc1.0.20200916142827-bd33bbf0497b+incompatible

replace github.com/terraform-providers/terraform-provider-nutanix => github.com/nutanix/terraform-provider-nutanix v1.2.2-0.20211029075448-e21f85ac2cf7

require (
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/docker/docker v1.13.1 // indirect
	github.com/docker/machine v0.16.2
	github.com/google/uuid v1.3.0
	github.com/sirupsen/logrus v1.8.1
	github.com/terraform-providers/terraform-provider-nutanix v1.1.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
