package types

type ProxmoxStatusResponse struct {
	Data []map[string]interface{}
}

type NodeWithContainers struct {
	Node       Node
	Containers []Container
}

type ProxmoxActionResponse struct {
	Success int    `json:"success"`
	Message string `json:"message"`
}