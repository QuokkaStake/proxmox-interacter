package types

type ProxmoxStatusResponse struct {
	Data []map[string]interface{}
}

type NodeWithAssets struct {
	Node       Node
	Containers []Container
	Storages   []Storage
}

type NodeWithContainers struct {
	Node       Node
	Containers []Container
}

type NodeWithStorages struct {
	Node     Node
	Storages []Storage
}

type ProxmoxActionResponse struct {
	Success int    `json:"success"`
	Message string `json:"message"`
}
