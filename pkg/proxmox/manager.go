package proxmox

import (
	"github.com/rs/zerolog"
	configPkg "main/pkg/config"
	"main/pkg/types"
)

type Manager struct {
	Logger  zerolog.Logger
	Clients []*Client
}

func NewManager(config *configPkg.Config, logger *zerolog.Logger) *Manager {
	clients := make([]*Client, len(config.Proxmox))

	for index, proxmoxConfig := range config.Proxmox {
		clients[index] = NewClient(proxmoxConfig, logger)
	}

	return &Manager{
		Logger:  logger.With().Str("component", "proxmox_manager").Logger(),
		Clients: clients,
	}
}

func (m *Manager) GetResources() ([]*types.ProxmoxStatusResponse, error) {
	responses := make([]*types.ProxmoxStatusResponse, 0)

	for _, client := range m.Clients {
		if response, err := client.GetResources(); err != nil {
			return responses, err
		} else {
			responses = append(responses, response)
		}
	}

	return responses, nil
}
