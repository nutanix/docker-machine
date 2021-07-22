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
5. Ability to specify the network(s) of the VM
6. Ability to specify the disk in the VM by image name
7. Ability to specify categories to applied to the VM ( flow, leap, ...)
8. Ability to add one additional disk by specifying disk-size and storage-container

Driver Args
-----------
|Arg                             |Description                                                              |Required          |
|--------------------------------|:------------------------------------------------------------------------|:-----------------|
| `--nutanix-endpoint`           |The hostname/ip-address of the Prism Central                             |yes               |
| `--nutanix-username`           |The username of the nutanix management account                           |yes               |
| `--nutanix-password`           |The password of the nutanix management account                           |yes               |
| `--nutanix-insecure`           |Set to true to force SSL insecure connection                             |no (default=false)|
| `--nutanix-cluster`            |The name of the cluster where deploy the VM (case sensitive)             |yes               |
| `--nutanix-vm-mem`             |The amount of RAM of the newly created VM                                |no (default=2G)   |
| `--nutanix-vm-cpus`            |The number of cpus in the newly created VM                               |no (default=2)    |
| `--nutanix-vm-cores`           |The number of cores per vCPU                                             |no (default=1)    |
| `--nutanix-vm-network`         |The network(s) to which the VM is attached to, support multiple network (separated by a comma)|yes               |
| `--nutanix-vm-image`           |The name of the Image to clone from                                      |yes               |
| `--nutanix-vm-categories`      |The name of the categories who will be applied to the newly created VM   |no                |
| `--nutanix-disk-size`          |The size of the additional disk to add to the VM                         |no                |
| `--nutanix-storage-container`  |The storage container UUID of the additional disk to add to the VM       |no                |
