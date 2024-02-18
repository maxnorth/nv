package resolver

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
)

func (r Resolver) ResolveEnvironment() (map[string]string, error) {
	resolvedValues := map[string]string{}
	for _, e := range os.Environ() {
		rawPair := strings.SplitN(e, "=", 2)
		name, value := rawPair[0], rawPair[1]

		// TODO: what are the error conditions vs ignore?
		if !strings.HasPrefix(value, "nv+") {
			continue
		}

		value, err := r.ResolveUrl(value)
		if err != nil {
			return nil, err
		}

		resolvedValues[name] = value
		os.Setenv(name, value)
	}

	return resolvedValues, nil
}

func (r *Resolver) ResolveUrl(rawUrl string) (string, error) {
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}

	schemes := strings.SplitN(parsedUrl.Scheme, "+", 2)
	if schemes[0] != "nv" {
		return "", errors.New("TODO")
	}

	resolverName := schemes[1]

	resolver, found := r.providers[resolverName]
	if !found {
		if !r.configFound {
			fmt.Fprintf(os.Stderr, "warning: no nv.yaml was found, use this to define custom resolvers\n")
		}
		return "", fmt.Errorf("no resolver found for '%s'", resolverName)
	}

	if _, loaded := r.loadedProviders[resolver]; !loaded {
		err := resolver.Load()
		if err != nil {
			return "", fmt.Errorf("failed to load provider '%s': %s", resolverName, err)
		}

		r.loadedProviders[resolver] = struct{}{}
	}

	value, err := resolver.GetValue(rawUrl)
	if err != nil {
		return "", err
	}

	return value, nil
}
