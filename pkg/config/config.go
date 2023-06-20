package config

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

type ProxmoxConfig struct {
	URL   string `default:"http://localhost:3000" toml:"url"`
	User  string `toml:"user"`
	Token string `toml:"token"`
}

func (c *Config) Validate() error {
	return nil
}
