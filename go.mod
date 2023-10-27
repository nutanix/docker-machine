module github.com/nutanix/docker-machine

go 1.18

replace github.com/docker/docker => github.com/docker/engine v17.12.0-ce-rc1.0.20200916142827-bd33bbf0497b+incompatible

require (
	github.com/docker/machine v0.16.2
	github.com/google/uuid v1.4.0
	github.com/nutanix-cloud-native/prism-go-client v0.2.0
	github.com/sirupsen/logrus v1.9.3
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/PaesslerAG/gval v1.2.0 // indirect
	github.com/PaesslerAG/jsonpath v0.1.1 // indirect
	github.com/docker/docker v20.10.17+incompatible // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/hashicorp/go-hclog v1.2.2 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/hashicorp/logutils v1.0.0 // indirect
	github.com/hashicorp/terraform-plugin-log v0.7.0 // indirect
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.20.0 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	golang.org/x/crypto v0.0.0-20220722155217-630584e8d5aa // indirect
	golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab // indirect
	golang.org/x/term v0.0.0-20220722155259-a9ba230a4035 // indirect
	gotest.tools v2.2.0+incompatible // indirect
)
