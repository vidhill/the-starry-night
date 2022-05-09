package domain

import (
	"fmt"

	"github.com/spf13/viper"
)

type ConfigViperRepository struct {
	Registry *viper.Viper
}

func (c ConfigViperRepository) GetString(s string) string {
	return c.Registry.GetString(s)
}

func (c ConfigViperRepository) GetBool(s string) bool {
	return c.Registry.GetBool(s)
}

func (c ConfigViperRepository) GetInt(s string) int {
	return c.Registry.GetInt(s)
}

func NewViperConfig() ConfigViperRepository {

	registry := viper.New()

	registry.SetConfigName("settings_private")
	registry.AddConfigPath("./")

	registry.ReadInConfig() // Find and read the (private) config file

	registry.SetConfigName("settings") // name of config file (without extension)

	err := registry.MergeInConfig() // merge settings.yaml into settings_private.yaml

	registry.AutomaticEnv() // override settings with env vars (if set)

	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	return ConfigViperRepository{
		Registry: registry,
	}
}
