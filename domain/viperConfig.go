package domain

type ViperConfig struct{}

func (v ViperConfig) GetString(s string) string {
	return "8080"
}

func (v ViperConfig) GetBool(s string) bool {
	return true
}

func NevViperConfig() ViperConfig {
	return ViperConfig{}
}
