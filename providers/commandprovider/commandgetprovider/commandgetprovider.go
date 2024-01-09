package commandgetprovider

import (
	"github.com/maxnorth/nv/providers"
	"github.com/maxnorth/nv/providers/commandprovider"
)

type provider struct {
	config commandprovider.Config
}

func New(config commandprovider.Config) providers.Provider {
	return &provider{
		config: config,
	}
}
