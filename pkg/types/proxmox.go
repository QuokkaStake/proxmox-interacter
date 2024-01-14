package types

type ProxmoxStatusResponse struct {
	Data []map[string]interface{}
}

type NodeWithAssets struct {
	Node       Node
	Containers []Container
	Storages   []Storage
}

type ProxmoxActionResponse struct {
	Success int    `json:"success"`
	Message string `json:"message"`
}

/*
LXC:
{
  "success": 1,
  "data": {
    "features": "nesting=1",
    "arch": "amd64",
    "hostname": "cosmos-testnet",
    "onboot": 1,
    "memory": 32768,
    "swap": 2048,
    "unprivileged": 1,
    "net0": "name=eth0,bridge=vmbr0,firewall=1,hwaddr=5A:7F:B5:30:70:A4,ip=dhcp,ip6=dhcp,type=veth",
    "rootfs": "disk1:vm-100-disk-0,size=400G",
    "ostype": "ubuntu",
    "cores": 4,
    "digest": "7710631d96bf0cbd91b2cd9e732fa080378dee37"
  }
}
*/

type ProxmoxLxcConfigResponse struct {
	Success int `json:"success"`
	Data    struct {
		Features     string `json:"features"`
		Arch         string `json:"arch"`
		Hostname     string `json:"hostname"`
		OnBoot       int    `json:"onboot"`
		Memory       uint64 `json:"memory"`
		Swap         uint64 `json:"swap"`
		Unprivileged int    `json:"unprivileged"`
		Net0         string `json:"net0"`
		RootFS       string `json:"rootfs"`
		OsType       string `json:"ostype"`
		Cores        int    `json:"cores"`
		Digest       string `json:"digest"`
	} `json:"data"`
}

/*
VM:

	{
	    "success": 1,
	    "data": {
	        "vmgenid": "43b7bc9f-7a22-4dc4-a76b-82c97197ceed",
	        "cores": 1,
	        "ostype": "l26",
	        "scsihw": "virtio-scsi-single",
	        "sockets": 1,
	        "name": "test-vm",
	        "digest": "3b53178a656bdb9326fd1c59d7fd3657ed1ca7dd",
	        "smbios1": "uuid=bfe2a387-775a-4ebe-958d-0a2b1c4c7bad",
	        "ide2": "none,media=cdrom",
	        "meta": "creation-qemu=8.1.2,ctime=1705059945",
	        "net0": "virtio=BC:24:11:83:B5:D8,bridge=vmbr0,firewall=1",
	        "scsi0": "disk3:vm-113-disk-0,iothread=1,size=10G",
	        "cpu": "x86-64-v2-AES",
	        "memory": "2048",
	        "boot": "order=scsi0;ide2;net0",
	        "numa": 0
	    }
	}
*/

type ProxmoxQemuConfigResponse struct {
	Success int `json:"success"`
	Data    struct {
		Name   string `json:"name"`
		Memory string `json:"memory"`
		Net0   string `json:"net0"`
		OsType string `json:"ostype"`
		Cores  int    `json:"cores"`
		Digest string `json:"digest"`
		OnBoot int    `json:"onboot"`
	} `json:"data"`
}
