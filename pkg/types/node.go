package types

import "fmt"

//{
//"maxmem": 67186364416,
//"id": "node/proxmox-1",
//"cpu": 0.0238454898575744,
//"type": "node",
//"node": "proxmox-1",
//"status": "online",
//"mem": 24367136768,
//"level": "",
//"maxdisk": 461901004800,
//"disk": 8873836544,
//"uptime": 1153492,
//"maxcpu": 20,
//"cgroup-mode": 2
//},

type Node struct {
	ID     string `json:"id"`
	Node   string `json:"node"`
	Uptime int64  `json:"uptime"`
	Status string `json:"status"`

	Memory int64   `type:"mem"`
	CPU    float64 `type:"cpu"`
	Disk   int64   `json:"disk"`

	MaxCPU    int64 `json:"maxcpu"`
	MaxMemory int64 `json:"maxmem"`
	MaxDisk   int64 `json:"maxdisk"`

	Link Link `json:"-"`
}

func (n Node) GetID() string   { return n.ID }
func (n Node) GetName() string { return n.Node }

func (n Node) IsRunning() bool {
	return n.Status == "online"
}

func (n Node) GetEmoji() string {
	if n.Status == "online" {
		return "🟢"
	}

	return "🔴"
}

func (n Node) GetCPUUsage() string {
	return fmt.Sprintf("%.2f%%", n.CPU*100)
}
