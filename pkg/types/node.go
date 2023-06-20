package types

import (
	"encoding/json"
	"fmt"
)

//{
//"maxmem": 67186339840,
//"uptime": 1035618,
//"status": "online",
//"node": "proxmox-3",
//"type": "node",
//"maxcpu": 20,
//"level": "",
//"mem": 24002428928,
//"cpu": 0.0697250418843008,
//"disk": 23254401024,
//"cgroup-mode": 2,
//"id": "node/proxmox-3",
//"maxdisk": 461916864512
//},

type Node struct {
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
