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

func (m *Manager) findContainerAndClientByName(query string) (*types.Container, *Client, error) {
	clusters, err := m.GetNodes()
	if err != nil {
		return nil, nil, err
	}

	container, clusterName, found := clusters.FindContainer(query)
	if !found {
		return nil, nil, fmt.Errorf("Container is not found!")
	}

	for _, client := range m.Clients {
		if client.Config.Name == clusterName {
			return container, client, nil
		}
	}

	return nil, nil, fmt.Errorf("Cluster is not found!")
}

func (m *Manager) StartContainer(containerName string) (*types.Container, error) {
	container, client, err := m.findContainerAndClientByName(containerName)
	if err != nil {
		return container, err
	}

	if container.IsRunning() {
		return container, fmt.Errorf("Container is already running!")
	}

	_, err = client.StartContainer(*container)
	if err != nil {
		return container, err
	}

	return container, nil
}

func (m *Manager) StopContainer(containerName string) (*types.Container, error) {
	container, client, err := m.findContainerAndClientByName(containerName)
	if err != nil {
		return container, err
	}

	if !container.IsRunning() {
		return container, fmt.Errorf("Container is not running!")
	}

	_, err = client.StopContainer(*container)
	if err != nil {
		return container, err
	}

	return container, nil
}

func (m *Manager) RestartContainer(containerName string) (*types.Container, error) {
	container, client, err := m.findContainerAndClientByName(containerName)
	if err != nil {
		return container, err
	}

	if !container.IsRunning() {
		return container, fmt.Errorf("Container is not running!")
	}

	_, err = client.RebootContainer(*container)
	if err != nil {
		return container, err
	}

	return container, nil
}
