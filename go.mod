module github.com/nutanix/docker-machine

go 1.21

replace github.com/docker/docker => github.com/docker/engine v17.12.0-ce-rc1.0.20200916142827-bd33bbf0497b+incompatible

require (
	github.com/docker/machine v0.16.2
	github.com/google/uuid v1.6.0
	github.com/sirupsen/logrus v1.9.3
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/docker/docker v20.10.17+incompatible // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	golang.org/x/crypto v0.23.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/term v0.20.0 // indirect
	gotest.tools v2.2.0+incompatible // indirect
)
