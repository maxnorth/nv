package resolver

import (
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

		if !strings.HasPrefix(value, "nv://") {
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

	provider, found := r.providers[parsedUrl.Host]
	if !found {
		// TODO: revisit - happy with this pattern?
		fmt.Printf("warning: no provider found for '%s'\n", parsedUrl.Host)
		return "", nil
	}

	if _, loaded := r.loadedProviders[provider]; !loaded {
		err := provider.Load()
		if err != nil {
			fmt.Printf("error: failed to load provider '%s': %s\n", parsedUrl.Host, err)
			os.Exit(1)
		}

		r.loadedProviders[provider] = struct{}{}
	}

	value, err := provider.GetValue(rawUrl)
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	return value, nil
}
