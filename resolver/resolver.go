package resolver

import "github.com/maxnorth/nv/providers"

type Resolver struct {
	providers map[string]providers.Provider
}
