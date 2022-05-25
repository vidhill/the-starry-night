package mocks

import "github.com/stretchr/testify/mock"

type Config struct {
	mock.Mock
}

func (mock *Config) GetBool(s string) bool {
	args := mock.Called(s)
	return args.Bool(0)
}
func (mock *Config) GetString(s string) string {
	args := mock.Called(s)
	return args.String(0)
}
func (mock *Config) GetInt(s string) int {
	args := mock.Called(s)
	return args.Int(0)
}

func NewMockConfig() Config {
	return Config{}
}
