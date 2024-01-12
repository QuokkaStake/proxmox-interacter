package proxmox

import (
	"fmt"
	"main/pkg/types"
	"sync"

	"github.com/rs/zerolog"
)

type Manager struct {
	Logger  zerolog.Logger
	Clients []*Client
}

func NewManager(config *types.Config, logger *zerolog.Logger) *Manager {
	clients := make([]*Client, len(config.Proxmox))

	for index, proxmoxConfig := range config.Proxmox {
		clients[index] = NewClient(proxmoxConfig, logger)
	}

	return &Manager{
		Logger:  logger.With().Str("component", "proxmox_manager").Logger(),
		Clients: clients,
	}
}

func (m *Manager) GetNodes() (types.ClusterInfos, error) {
	responses := make(types.ClusterInfos, len(m.Clients))

	var mutex sync.Mutex
	var wg sync.WaitGroup

	for index, client := range m.Clients {
		wg.Add(1)
		go func(index int, client *Client) {
			defer wg.Done()

			response, err := client.GetNodesWithAssets()
			clusterWithError := types.ClusterInfo{Name: client.Config.Name}

			if err != nil {
				clusterWithError.Error = err
			} else {
				clusterWithError.Nodes = response
			}

			mutex.Lock()
			responses[index] = clusterWithError
			mutex.Unlock()
		}(index, client)
	}

	wg.Wait()

	return responses, nil
}
func (m *Manager) findClientByName(clusterName string) (*Client, bool) {
	for _, client := range m.Clients {
		if client.Config.Name == clusterName {
			return client, true
		}
	}

	return nil, false
}

func (m *Manager) findContainerAndClientByName(query string) (*types.Container, *Client, error) {
	clusters, err := m.GetNodes()
	if err != nil {
		return nil, nil, err
	}

	container, clusterName, err := clusters.FindContainer(query)
	if err != nil {
		return nil, nil, err
	}

	if client, found := m.findClientByName(clusterName); found {
		return container, client, nil
	}

	return nil, nil, fmt.Errorf("Cluster is not found!")
}

func (m *Manager) StartContainer(containerName string) error {
	container, client, err := m.findContainerAndClientByName(containerName)
	if err != nil {
		return err
	}

	_, err = client.StartContainer(*container)
	return err
}

func (m *Manager) StopContainer(containerName string) error {
	container, client, err := m.findContainerAndClientByName(containerName)
	if err != nil {
		return err
	}

	_, err = client.StopContainer(*container)
	return err
}

func (m *Manager) RestartContainer(containerName string) error {
	container, client, err := m.findContainerAndClientByName(containerName)
	if err != nil {
		return err
	}

	_, err = client.RebootContainer(*container)
	return err
}

func (m *Manager) GetContainerConfig(container types.Container, clusterName string) (*types.ContainerConfig, error) {
	client, found := m.findClientByName(clusterName)
	if !found {
		return nil, fmt.Errorf("Cluster is not found!")
	}

	switch container.Type {
	case "lxc":
		if config, err := client.GetLxcContainerConfig(container); err != nil {
			return nil, err
		} else {
			return client.ParseLxcContainerConfig(config)
		}
	case "qemu":
		if config, err := client.GetQemuContainerConfig(container); err != nil {
			return nil, err
		} else {
			return client.ParseQemuContainerConfig(config)
		}
	}

	return nil, fmt.Errorf("Unsupported container type: %s", container.Type)
}
