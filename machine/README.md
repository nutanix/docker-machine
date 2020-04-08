# Nutanix Docker Machine Driver

This repository contains the docker machine driver for nutanix.

The driver can create VMs in a nutanix managed cluster, and setup docker on it. It uses the MGMT API to achieve this goal.

Build Instructions
--------------------

make sure the vendored (code from the vendor directory) version of docker-machine code is in $GOPATH

`go build -o docker-machine-driver-nutanix main.go` 

Features
---------

1. Ability to select VM's Main Memory in Megabytes
2. Ability to select VM's vCPU count
3. Ability to set a custom name for the newly created VM
4. Ability to set the number of cores per vCPU
5. Ability to specify the network of the VM
6. Ability to specify the disk in the VM by image name

Driver Args
-----------
|Arg                             |Description                                                              |Required          |
|--------------------------------|:------------------------------------------------------------------------|:-----------------|
| `--nutanix-username`           |The username of the nutanix management account                           |yes               |
| `--nutanix-password`           |The password of the nutanix management account                           |yes               |
| `--nutanix-endpoint`           |The hostname/ip-address of the management API server of the cluster      |yes               |
| `--nutanix-vm-mem`             |The amount of RAM of the newly created VM                                |no (default=1G)   |
| `--nutanix-vm-cpus`            |The number of cpus in the newly created VM                               |no (default=1)    |
| `--nutanix-vm-cores`           |The number of cores per vCPU                                             |no (default=1)    |
| `--nutanix-vm-network`         |The network to which the vNIC of the VM is attached to                   |yes               |
| `--nutanix-vm-image`           |The name of the Image to clone from                                      |yes               |

