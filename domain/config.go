package domain

type ConfigProvider interface {
	GetString(string) string
	GetBool(string) bool
	GetInt(string) int
}
