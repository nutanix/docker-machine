package driver

import "gopkg.in/yaml.v3"

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
