package types

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
