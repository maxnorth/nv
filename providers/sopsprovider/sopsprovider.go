package sopsprovider

import "github.com/maxnorth/nv/providers"

type Config struct {
	File string
}

type provider struct {
	config   Config
	jsonData []byte
}

func New(config Config) providers.Provider {
	return &provider{
		config: config,
	}
}
