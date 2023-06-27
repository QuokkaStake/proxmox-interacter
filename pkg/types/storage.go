package types

// {
//	"maxmem": 17179869184,
//	"name": "neutron-validator",
//	"id": "lxc/108",
//	"netout": 252286415555,
//	"disk": 178980913152,
//	"node": "proxmox-2",
//	"uptime": 980343,
//	"mem": 7803162624,
//	"type": "lxc",
//	"template": 0,
//	"status": "running",
//	"diskread": 1721992290304,
//	"netin": 358679496669,
//	"maxdisk": 421606629376,
//	"diskwrite": 10378809516032,
//	"vmid": 108,
//	"maxcpu": 2,
//	"cpu": 0.0474592726849931
// },

type Storage struct {
	ID         string `json:"id"`
	Node       string `json:"node"`
	Storage    string `json:"storage"`
	Status     string `json:"status"`
	PluginType string `json:"plugintype"`

	MaxDisk int64 `json:"maxdisk"`
	Disk    int64 `json:"disk"`

	Link Link `json:"-"`
}

func (s Storage) GetID() string   { return s.ID }
func (s Storage) GetName() string { return s.Storage }

func (s Storage) GetEmoji() string {
	if s.Status == "available" {
		return "🟢"
	}

	return "⚪"
}
