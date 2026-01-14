package config

import "fmt"

func (cfg ServerConfig) PortStr() string {
	return fmt.Sprintf(":%d", cfg.Port)
}
