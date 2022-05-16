package stubrepository

type StubConfig struct{}

func (mock *StubConfig) GetBool(s string) bool {
	return false
}
func (mock *StubConfig) GetString(s string) string {
	return ""
}
func (mock *StubConfig) GetInt(s string) int {
	return 0
}

func NewStubConfig() StubConfig {
	return StubConfig{}
}
