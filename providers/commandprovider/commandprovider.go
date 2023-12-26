package commandprovider

import "github.com/maxnorth/nv/providers"

type Config struct {
	Command string
}

type provider struct {
	config Config
	values map[string]string
}

func New(config Config) providers.Provider {
	return &provider{
		config: config,
	}
}
