package types

import "fmt"

//{
//"maxmem": 17179869184,
//"name": "neutron-validator",
//"id": "lxc/108",
//"netout": 252286415555,
//"disk": 178980913152,
//"node": "proxmox-2",
//"uptime": 980343,
//"mem": 7803162624,
//"type": "lxc",
//"template": 0,
//"status": "running",
//"diskread": 1721992290304,
//"netin": 358679496669,
//"maxdisk": 421606629376,
//"diskwrite": 10378809516032,
//"vmid": 108,
//"maxcpu": 2,
//"cpu": 0.0474592726849931
//},

type Container struct {
	ID     string `json:"id"`
	VMID   int64  `json:"vmid"`
	Node   string `json:"node"`
	Type   string `json:"type"`
	Name   string `json:"name"`
	Uptime int64  `json:"uptime"`
	Status string `json:"status"`

	MaxCPU    int64 `json:"maxcpu"`
	MaxMemory int64 `json:"maxmem"`
	MaxDisk   int64 `json:"maxdisk"`

	NetIn     int64   `json:"netin"`
	NetOut    int64   `json:"netout"`
	DiskRead  int64   `json:"diskread"`
	DiskWrite int64   `json:"diskwrite"`
	Disk      int64   `json:"disk"`
	Memory    int64   `json:"mem"`
	CPU       float64 `json:"cpu"`

	Link Link `json:"-"`
}

func (c Container) GetID() string   { return c.ID }
func (c Container) GetName() string { return c.Name }

func (c Container) GetCPUUsage() string {
	return fmt.Sprintf("%.2f%%", c.CPU*100)
}

type NodeWithContainers struct {
	Node       Node
	Containers []Container
}

func (c Container) GetEmoji() string {
	if c.Status == "running" {
		return "🟢"
	}

	return "⚪"
}

func (c Container) IsRunning() bool {
	return c.Status == "running"
}

func (c Container) GetTypeEmoji() string {
	if c.Type == "lxc" {
		return "📦"
	}

	return "🖥️"
}

func (c Container) GetType() string {
	if c.Type == "lxc" {
		return "LXC container"
	}

	return "Virtual machine"
}

type Node struct {
	ID        string `json:"id"`
	Name      string `json:"node"`
	Uptime    int64  `json:"uptime"`
	Status    string `json:"status"`
	MaxCPU    int64  `json:"maxcpu"`
	MaxMemory int64  `json:"maxmem"`
	MaxDisk   int64  `json:"maxdisk"`

	Link Link `json:"-"`
}

func (n Node) GetID() string   { return n.ID }
func (n Node) GetName() string { return n.Name }

func (n Node) GetEmoji() string {
	if n.Status == "online" {
		return "🟢"
	}

	return "🔴"
}

type Resource interface {
	GetID() string
	GetName() string
}
