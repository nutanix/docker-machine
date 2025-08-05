package driver

import (
	"context"
	"fmt"
	"regexp"

	log "github.com/sirupsen/logrus"

	v3 "github.com/nutanix-cloud-native/prism-go-client/v3"
	"gopkg.in/yaml.v3"
)

func isUUID(uuid string) bool {
	uuidPattern := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	return uuidPattern.MatchString(uuid)
}

func iterateNode(node *yaml.Node, identifier string) *yaml.Node {
	returnNode := false
	for _, n := range node.Content {
		if n.Value == identifier {
			returnNode = true
			continue
		}
		if returnNode {
			return n
		}
		if len(n.Content) > 0 {
			ac_node := iterateNode(n, identifier)
			if ac_node != nil {
				return ac_node
			}
		}
	}
	return nil
}

// deleteAllContents will remove all the contents of a node
// Mark sure to pass the correct node in otherwise bad things will happen
// func deleteAllContents(node *yaml.Node) {
// 	node.Content = []*yaml.Node{}
// }

// buildStringNodes builds Nodes for a single key: value instance
func buildStringNodes(key, value, comment string) []*yaml.Node {
	keyNode := &yaml.Node{
		Kind:        yaml.ScalarNode,
		Tag:         "!!str",
		Value:       key,
		HeadComment: comment,
	}
	valueNode := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!str",
		Value: value,
	}
	return []*yaml.Node{keyNode, valueNode}
}

func buildScalarNodes(key string) []*yaml.Node {
	keyNode := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!str",
		Value: key,
	}
	return []*yaml.Node{keyNode}
}

// buildMapNodes builds Nodes for a key: map instance
// func buildMapNodes(key string) (*yaml.Node, *yaml.Node) {
// 	n1, n2 := &yaml.Node{
// 		Kind:  yaml.ScalarNode,
// 		Tag:   "!!str",
// 		Value: key,
// 	}, &yaml.Node{Kind: yaml.MappingNode,
// 		Tag: "!!map",
// 	}
// 	return n1, n2
// }

// buildSeqNodes builds Nodes for a key: map instance
// func buildSeqNodes(key string) (*yaml.Node, *yaml.Node) {
// 	n1, n2 := &yaml.Node{
// 		Kind:  yaml.ScalarNode,
// 		Tag:   "!!str",
// 		Value: key,
// 	}, &yaml.Node{Kind: yaml.SequenceNode,
// 		Tag: "!!seq",
// 	}
// 	return n1, n2
// }

// GetGPUList retrieves a list of GPUs from the Nutanix Prism Element based on the provided GPU names and PE UUID.
// It returns a slice of VMGpu pointers or an error if any issues occur during the retrieval
func GetGPUList(ctx context.Context, conn *v3.Client, gpus []string, peUUID string) ([]*v3.VMGpu, error) {
	resultGPUs := make([]*v3.VMGpu, 0)
	for _, gpu := range gpus {
		foundGPU, err := GetGPU(ctx, conn, peUUID, gpu)
		if err != nil {
			return nil, err
		}
		resultGPUs = append(resultGPUs, foundGPU)
	}
	return resultGPUs, nil
}

// GetGPU retrieves a specific GPU from the Nutanix Prism Element based on the provided GPU name and PE UUID.
// It returns a VMGpu pointer or an error if the GPU is not found or if any issues occur during the retrieval.
func GetGPU(ctx context.Context, conn *v3.Client, peUUID, gpu string) (*v3.VMGpu, error) {
	if gpu == "" {
		return nil, fmt.Errorf("gpu name must be passed in order to retrieve the GPU")
	}

	log.Infof("Searching GPU %s in Prism Element with UUID %s", gpu, peUUID)
	allGPUs, err := GetGPUsForPE(ctx, conn, peUUID)
	if err != nil {
		return nil, err
	}
	if len(allGPUs) == 0 {
		return nil, fmt.Errorf("no available GPUs found in Prism Element cluster with UUID %s", peUUID)
	}
	for _, peGPU := range allGPUs {
		if peGPU.Status != "UNUSED" {
			continue
		}
		if peGPU.Name == gpu {
			log.Infof("GPU %s found with ID %d in Prism Element", peGPU.Name, *peGPU.DeviceID)
			return &v3.VMGpu{
				DeviceID: peGPU.DeviceID,
				Mode:     &peGPU.Mode,
				Vendor:   &peGPU.Vendor,
			}, nil
		}
	}
	return nil, fmt.Errorf("no available GPU found in Prism Element that matches required GPU name: %s", gpu)

}

// GetGPUsForPE retrieves all GPUs associated with a specific Prism Element (PE) UUID.
// It returns a slice of GPU pointers or an error if any issues occur during the retrieval.
func GetGPUsForPE(ctx context.Context, conn *v3.Client, peUUID string) ([]*v3.GPU, error) {
	gpus := make([]*v3.GPU, 0)
	hosts, err := conn.V3.ListAllHost(ctx)
	if err != nil {
		return gpus, err
	}

	for _, host := range hosts.Entities {
		if host == nil ||
			host.Status == nil ||
			host.Status.ClusterReference == nil ||
			host.Status.Resources == nil ||
			len(host.Status.Resources.GPUList) == 0 ||
			host.Status.ClusterReference.UUID != peUUID {
			continue
		}

		for _, peGpu := range host.Status.Resources.GPUList {
			if peGpu == nil {
				continue
			}
			gpus = append(gpus, peGpu)
		}
	}
	return gpus, nil
}
