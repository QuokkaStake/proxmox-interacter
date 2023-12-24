package types

import (
	"strconv"
)

type ClusterInfos []ClusterInfo

type ClusterInfo struct {
	Name  string
	Nodes []NodeWithAssets
	Error error
}

func (c ClusterInfos) FindNode(query string) (*Node, bool) {
	for _, cluster := range c {
		if cluster.Error != nil {
			continue
		}

		for _, node := range cluster.Nodes {
			if node.Node.Node == query || node.Node.ID == query {
				return &node.Node, true
			}
		}
	}

	return nil, false
}

func (c ClusterInfos) FindContainer(query string) (*Container, bool) {
	for _, cluster := range c {
		if cluster.Error != nil {
			continue
		}

		for _, node := range cluster.Nodes {
			for _, container := range node.Containers {
				// containers IDs are like lxc/XXXX or qemu/XXXX
				if strconv.FormatInt(container.VMID, 10) == query ||
					container.ID == query ||
					container.Name == query {
					return &container, true
				}
			}
		}
	}

	return nil, false
}
