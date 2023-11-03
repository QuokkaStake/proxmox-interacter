package proxmox

import (
	"fmt"
	"main/pkg/types"
	"main/pkg/utils"
	"strconv"

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

func (m *Manager) GetNodes() ([]types.Node, error) {
	responses := make([]types.Node, 0)

	for _, client := range m.Clients {
		if response, err := client.GetNodes(); err != nil {
			return responses, err
		} else {
			responses = append(responses, response...)
		}
	}

	return responses, nil
}

func (m *Manager) GetContainers() ([]types.Container, error) {
	responses := make([]types.Container, 0)

	for _, client := range m.Clients {
		if response, err := client.GetContainers(); err != nil {
			return responses, err
		} else {
			responses = append(responses, response...)
		}
	}

	return responses, nil
}

func (m *Manager) GetNodesWithContainers() ([]types.NodeWithContainers, error) {
	responses := make([]types.NodeWithContainers, 0)

	for _, client := range m.Clients {
		if response, err := client.GetNodesWithContainers(); err != nil {
			return responses, err
		} else {
			responses = append(responses, response...)
		}
	}

	return responses, nil
}

func (m *Manager) FindNodeWithClient(containerName string) (*types.Container, *Client, error) {
	for _, client := range m.Clients {
		containers, err := client.GetContainers()
		if err != nil {
			return nil, nil, err
		}

		container, found := utils.Find(containers, func(c types.Container) bool {
			return c.ID == containerName ||
				strconv.FormatInt(c.VMID, 10) == containerName ||
				c.Name == containerName
		})

		if found {
			return container, client, nil
		}
	}

	return nil, nil, fmt.Errorf("Container is not found")
}

func (m *Manager) StartContainer(containerName string) (*types.Container, error) {
	container, client, err := m.FindNodeWithClient(containerName)
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
	container, client, err := m.FindNodeWithClient(containerName)
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
	container, client, err := m.FindNodeWithClient(containerName)
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

func (m *Manager) GetNodesWithStorages() ([]types.NodeWithStorages, error) {
	responses := make([]types.NodeWithStorages, 0)

	for _, client := range m.Clients {
		if response, err := client.GetNodesWithStorages(); err != nil {
			return responses, err
		} else {
			responses = append(responses, response...)
		}
	}

	return responses, nil
}
