package commandlistprovider

import (
	"github.com/maxnorth/nv/providers"
	"github.com/maxnorth/nv/providers/commandprovider"
)

type provider struct {
	config commandprovider.Config
	values map[string]string
}

func New(config commandprovider.Config) providers.Provider {
	return &provider{
		config: config,
	}
}
