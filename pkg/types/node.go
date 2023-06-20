package types

import (
	"encoding/json"
	"fmt"
)

type Node struct {
	ID        string `json:"id"`
	Name      string `json:"node"`
	Uptime    int64  `json:"uptime"`
	Status    string `json:"status"`
	MaxCPU    int64  `json:"maxcpu"`
	MaxMemory int64  `json:"maxmem"`
	MaxDisk   int64  `json:"maxdisk"`
}

func ParseNodesFromResponse(response *ProxmoxStatusResponse) ([]Node, error) {
	nodes := make([]Node, 0)

	for _, rawData := range response.Data {
		resourceType, ok := rawData["type"]
		if !ok {
			return nodes, fmt.Errorf("resource type is not present")
		}

		if resourceType != "node" {
			continue
		}

		rawBytes, err := json.Marshal(rawData)
		if err != nil {
			return nodes, err
		}

		node := Node{}
		if err := json.Unmarshal(rawBytes, &node); err != nil {
			return nodes, err
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}

func (n Node) GetEmoji() string {
	if n.Status == "online" {
		return "🟢"
	}

	return "🔴"
}

func (n Node) GetLink(config ProxmoxConfig) Link {
	return config.GetResourceLink(n.ID, n.Name)
}

type NodeWithLink struct {
	Node Node
	Link Link
}
