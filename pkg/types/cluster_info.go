package types

type ClusterInfo struct {
	Name  string
	Nodes []NodeWithAssets
	Error error
}
