package providers

type Provider interface {
	Load() error
	GetValue(key string) (string, error)
}
