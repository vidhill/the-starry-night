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

func NewViperConfig() ConfigViperRepository {

	registry := viper.New()

	registry.SetConfigName("settings") // name of config file (without extension)
	registry.AddConfigPath("./")

	err := registry.ReadInConfig() // Find and read the config file
	registry.AutomaticEnv()        // override settings with env vars (if set)

	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	return ConfigViperRepository{
		Registry: registry,
	}
}
