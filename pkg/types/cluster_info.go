package types

type ClusterInfo struct {
	Name  string
	Nodes []Node
	Error error
}
