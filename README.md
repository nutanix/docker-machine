# Nutanix Rancher Node Driver

This repository contains the Rancher Node Driver for Nutanix. Nutanix Node driver are used to provision hosts on Nutanix Enterprise Cloud, which Rancher uses to launch and manage Kubernetes clusters.


---

[![Go Report Card](https://goreportcard.com/badge/github.com/nutanix/docker-machine)](https://goreportcard.com/report/github.com/nutanix/docker-machine)
![CI](https://github.com/nutanix/docker-machine/actions/workflows/integration.yml/badge.svg)
![Release](https://github.com/nutanix/docker-machine/actions/workflows/release.yml/badge.svg)

[![release](https://img.shields.io/github/release-pre/nutanix/docker-machine.svg)](https://github.com/nutanix/docker-machine/releases)
[![License](https://img.shields.io/badge/License-MPL%202.0-blue.svg)](https://github.com/nutanix/docker-machine/blob/master/LICENSE)
![Proudly written in Golang](https://img.shields.io/badge/written%20in-Golang-92d1e7.svg)
[![Releases](https://img.shields.io/github/downloads/nutanix/docker-machine/total.svg)](https://github.com/nutanix/docker-machine/releases)

---

## Features


- Configure Prism Central and corresponding user to talk to Nutanix platform
- Define target cluster to deploy VM
- Ability to set a custom name for the newly created VM
- Ability to select VM's Main Memory in Megabytes
- Ability to select VM's vCPU count
- Ability to set the number of cores per vCPU
- Ability to specify the network(s) of the VM (Classic or VPC)
- Ability to specify the template disk in the VM by image name and modify his size (increase only)
- Ability to specify categories to applied to the VM ( flow, leap, ...)
- Ability to add one additional disk by specifying disk-size and storage-container
- Enable passthrough the host's CPU features to the newly created VM
- Define a Cloud-init user-data to send to the newly created VM
- Project support
- Serial Port support
- Boot type selection : Legacy or UEFI 
- GPU support
- Prism Central Service Accounts support


## Installation


If you want to use Nutanix Node Driver, you need add it in order to start using them to create node templates and eventually node pools for your Kubernetes cluster.

1. From the Home view, choose *Cluster Management* > *Drivers* in the navigation bar. From the Drivers page, select the *Node Drivers* tab.
2. Click *Add Node Driver*.
3. Complete the Add Node Driver form. Then click Create.

    - *Download URL*: `https://github.com/nutanix/docker-machine/releases/download/v3.7.0/docker-machine-driver-nutanix`  
    - *Custom UI URL*: `https://nutanix.github.io/rancher-ui-driver/v3.7.0/component.js`
    - *Checksum*: `2f70c4bdccd3c5e68bd8c32aadb5b525275a3cda5799f29736f37bdd168caa94`  
    - *Whitelist Domains*: `nutanix.github.io`  
      
    
    *whitelist is mandatory and need to be changed if you relocate the UI driver*

<img width="956" height="461" alt="image" src="https://github.com/user-attachments/assets/1237ca91-e866-4906-aa1e-e9dfcdf2a58c" />





4. Wait for the driver to become "Active"
5. Go to *RKE1 Configuration > Node Templates*, your can create a Nutanix Template and custom UI should show up.

![image](https://github.com/nutanix/docker-machine/assets/180613/8c56a022-ad6b-406b-80e6-10c5673c0d9e)



## Driver Args

| Arg                          | Description                                                                                   | Required | Default |
|------------------------------|:----------------------------------------------------------------------------------------------|:---------|---------|
| `nutanix-endpoint`           | The hostname/ip-address of the Prism Central                                                  | yes      |         |
| `nutanix-port`               | The port to connect to Prism Central                                                          | no       | 9440    |
| `nutanix-username`           | The username of the nutanix management account                                                | yes      |         |
| `nutanix-password`           | The password of the nutanix management account                                                | yes      |         |
| `nutanix-insecure`           | Set to true to force SSL insecure connection                                                  | no       | false   |
| `nutanix-cluster`            | The name of the cluster where deploy the VM (case sensitive)                                  | yes      |         |
| `nutanix-boot-type`          | The boot type of the VM (legacy or uefi)                                                      | no       | legacy  |
| `nutanix-vm-mem`             | The amount of RAM of the newly created VM (MB)                                                | no       | 2 GB    |
| `nutanix-vm-cpus`            | The number of cpus in the newly created VM (core)                                             | no       | 2       |
| `nutanix-vm-cores`           | The number of cores per vCPU                                                                  | no       | 1       |
| `nutanix-vm-network`         | The network(s) to which the VM is attached to ( name or UUID )                                | yes      |         |
| `nutanix-vm-image`           | The name of the Disk Image template we use for the newly created VM (must support cloud-init) | yes      |         |
| `nutanix-vm-image-size`      | The new size of the Image we use as a template (in GiB)                                       | no       |         |
| `nutanix-vm-categories`      | The name of the categories who will be applied to the newly created VM                        | no       |         |
| `nutanix-vm-gpu`             | The name of GPU to attach to the newly created VM                                             | no       |         |
| `nutanix-project`            | The name of the project where deploy the VM (default if empty)                                | no       | default |
| `nutanix-disk-size`          | The size of the additional disk to add to the VM (in GiB)                                     | no       |         |
| `nutanix-storage-container`  | The storage container UUID of the additional disk to add to the VM                            | no       |         |
| `nutanix-cloud-init`         | Cloud-init to provide to the VM (will be patched with rancher root user)                      | no       |         |
| `nutanix-vm-cpu-passthrough` | Enable passthrough the host's CPU features to the newly created VM                            | no       | false   |
| `nutanix-vm-serial-port`     | Attach a serial port to the newly created VM                                                  | no       | false   |



## Project support

Starting `v3.3.0` the Rancher Node driver implements Nutanix Project support. The  prerequisite needed to be able to use the Rancher Node Driver is the following:
- Target cluster and network available in the Project
- Role with the following recommended permission:
  - VM Full Access
  - Cluster View Access
  - Image View Only

## Service Accounts support

Starting version 3.8.0, the Rancher Node Driver support Prism Central Service Accounts. 
To use a Service Account, you need to provide `X-ntnx-api-key` as the user name and the corresponding API Key as the password.

## Development

### Build Instructions

build linux/amd64 binary => `make` 
build local binary => `make local`

## History

* v1 is the original Nutanix docker machine driver that connect to Prism Element
* v2.x add Rancher 2.0 support
* v3.x is a rewrite of the driver that connect to Prism Central

## Support
This code is developed in the open with input from the community through issues and PRs. A Nutanix engineering team serves as the maintainer. Documentation is available in the project repository. Issues and enhancement requests can be submitted in the [Issues tab of this repository](../../issues). Please search for and review the existing open issues before submitting a new issue.

## License
Copyright 2022-2023 Nutanix, Inc.

The project is released under Mozilla Public License Version 2.0.
