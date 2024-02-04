package commandprovider

import (
	"github.com/maxnorth/nv/internal/providers"
)

type provider struct {
	config Config
}

type Config struct {
	Command string
	Mode    string
	Output  string
}

func New(config Config) providers.Provider {
	return &provider{
		config: config,
	}
}
