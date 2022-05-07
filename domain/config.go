package domain

type ConfigRepository interface {
	GetString(string) string
	GetBool(string) bool
}
