package resolver

import "github.com/maxnorth/nv/internal/providers"

type Resolver struct {
	providers       map[string]providers.Provider
	loadedProviders map[providers.Provider]struct{}
	LoadedVars      []string
}
