package hcvaultprovider

import (
	"net/http"

	"github.com/maxnorth/nv/providers"
)

type Config struct {
	Command []string
}

type provider struct {
	config Config
	client *http.Client
	values map[string]string
}

func New(config Config) providers.Provider {
	return &provider{
		config: config,
	}
}
