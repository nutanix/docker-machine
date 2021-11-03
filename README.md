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

Features
---------

1. Ability to select VM's Main Memory in Megabytes
2. Ability to select VM's vCPU count
3. Ability to set a custom name for the newly created VM
4. Ability to set the number of cores per vCPU
5. Ability to specify the network(s) of the VM
6. Ability to specify the template disk in the VM by image name and modify his size (increase only)
7. Ability to specify categories to applied to the VM ( flow, leap, ...)
8. Ability to add one additional disk by specifying disk-size and storage-container
9. Enable passthrough the host's CPU features to the newly created VM


Installation
--------------------

If you want to use Nutanix Node Driver, you need add it in order to start using them to create node templates and eventually node pools for your Kubernetes cluster.

1. From the Home view, choose *Cluster Management* > *Drivers* in the navigation bar. From the Drivers page, select the *Node Drivers* tab.
2. Click *Add Node Driver*.
3. Complete the Add Node Driver form. Then click Create.

    - *Download URL*: `https://github.com/nutanix/docker-machine/releases/download/v3.0.0/docker-machine-driver-nutanix_v3.0.0_linux`  
    - *Custom UI URL*: `https://nutanix.github.io/rancher-ui-driver/v3.0.0/component.js`  
    - *Whitelist Domains*: `nutanix.github.io`  
      
    *whitelist is mandatory and need to be changed if you relocate the UI driver*

![image](https://user-images.githubusercontent.com/180613/139593826-9d48bc40-29c0-42cb-8122-0e95304eeac8.png)

4. Wait for the driver to become "Active"
5. Go to *RKE1 Configuration > Node Templates*, your can create a Nutanix Template and custom UI should show up.

![image](https://user-images.githubusercontent.com/180613/139594240-db4f375f-5918-4918-b1be-4aa8e4232f0f.png)



Driver Args
-----------
|Arg                           |Description                                                              |Required          |Default |
|------------------------------|:------------------------------------------------------------------------|:-----------------|--------|
| `nutanix-endpoint`           |The hostname/ip-address of the Prism Central                             |yes               ||
| `nutanix-port`               |The port to connect to Prism Central                                     |no                |9440
| `nutanix-username`           |The username of the nutanix management account                           |yes               ||
| `nutanix-password`           |The password of the nutanix management account                           |yes               ||
| `nutanix-insecure`           |Set to true to force SSL insecure connection                             |no                |false|
| `nutanix-cluster`            |The name of the cluster where deploy the VM (case sensitive)             |yes               ||
| `nutanix-vm-mem`             |The amount of RAM of the newly created VM (MB)                           |no                | 2 GB|
| `nutanix-vm-cpus`            |The number of cpus in the newly created VM (core)                        |no                | 2|
| `nutanix-vm-cores`           |The number of cores per vCPU                                             |no                | 1|
| `nutanix-vm-network`         |The network(s) to which the VM is attached to                            |yes               ||
| `nutanix-vm-image`           |The name of the Image template we use for the newly created VM (must support cloud-init)|yes               ||
| `nutanix-vm-image-size`      |The new size of the Image we use as a template (in GiB)                  |no                ||
| `nutanix-vm-categories`      |The name of the categories who will be applied to the newly created VM   |no                ||
| `nutanix-disk-size`          |The size of the additional disk to add to the VM (in GiB)                |no                ||
| `nutanix-storage-container`  |The storage container UUID of the additional disk to add to the VM       |no                ||
| `nutanix-cloud-init`         |Cloud-init to provide to the VM (will be patched with rancher root user) |no                ||
| `nutanix-vm-cpu-passthrough` |Enable passthrough the host's CPU features to the newly created VM       |no                |false|

Build Instructions
--------------------

build linux/amd64 binary => `make`  
build local binary => `make local`
## History

* v1 is the original Nutanix docker machine driver that connect to Prism Element
* v2.x add Rancher 2.0 support
* v3.x is a rewrite of the driver that connect to Prism Central

