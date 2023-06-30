package types

import (
	"fmt"
	"net/url"
)

type Config struct {
	Log      LogConfig       `toml:"log"`
	Telegram TelegramConfig  `toml:"telegram"`
	Proxmox  []ProxmoxConfig `toml:"proxmox"`
}

type LogConfig struct {
	LogLevel   string `default:"info"  toml:"level"`
	JSONOutput bool   `default:"false" toml:"json"`
}

type TelegramConfig struct {
	Chat   int     `toml:"chat"`
	Token  string  `toml:"token"`
	Admins []int64 `toml:"admins"`
}

func (c TelegramConfig) Validate() error {
	if c.Token == "" {
		return fmt.Errorf("telegram token not specified")
	}

	return nil
}

type ProxmoxConfig struct {
	URL         string `default:"http://localhost:8006" toml:"url"`
	ExternalURL string `toml:"external-url"`
	User        string `toml:"user"`
	Token       string `toml:"token"`
}

func (c ProxmoxConfig) Host() string {
	if c.ExternalURL != "" {
		return c.ExternalURL
	}

	return c.URL
}

func (c ProxmoxConfig) Validate() error {
	if c.URL == "" {
		return fmt.Errorf("proxmox url not specified")
	}

	if c.User == "" {
		return fmt.Errorf("proxmox user not specified")
	}

	if c.Token == "" {
		return fmt.Errorf("proxmox token not specified")
	}

	return nil
}

func (c ProxmoxConfig) GetResourceLink(resource Resource) Link {
	return Link{
		Name: resource.GetName(),
		Href: fmt.Sprintf(
			"%s/#v1:0:=%s",
			c.Host(),
			url.QueryEscape(resource.GetID()),
		),
	}
}

func (c *Config) Validate() error {
	if err := c.Telegram.Validate(); err != nil {
		return fmt.Errorf("error in telegram config: %s", err)
	}

	for index, proxmoxConfig := range c.Proxmox {
		if err := proxmoxConfig.Validate(); err != nil {
			return fmt.Errorf("error in proxmox config #%d: %s", index, err)
		}
	}

	return nil
}
