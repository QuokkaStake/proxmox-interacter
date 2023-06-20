package proxmox

import (
	"github.com/rs/zerolog"
	"main/pkg/types"
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

func (m *Manager) GetNodes() ([]types.NodeWithLink, error) {
	responses := make([]types.NodeWithLink, 0)

	for _, client := range m.Clients {
		if response, err := client.GetNodes(); err != nil {
			return responses, err
		} else {
			responses = append(responses, response...)
		}
	}

	return responses, nil
}
