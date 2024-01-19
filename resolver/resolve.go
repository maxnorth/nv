package resolver

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type ValueConfig struct {
	ProviderAlias string
	ValueUrl      *url.URL
	ValueUrlRaw   string
	EnvKey        string
}

func (r Resolver) Resolve() map[string]string {
	err := godotenv.Load(".env", ".env.local")
	if err != nil {
		panic(err)
	}

	for _, provider := range r.providers {
		err := provider.Load()
		if err != nil {
			panic(err)
		}
	}

	// load providers eagerly - start dumb
	// check environment for providers and keys
	valueConfigs := []ValueConfig{}
	for _, e := range os.Environ() {
		if valueConfig, found := parseValueConfig(e); found {
			valueConfigs = append(valueConfigs, valueConfig)
		}
	}

	values := map[string]string{}
	for _, valueConfig := range valueConfigs {
		provider, found := r.providers[valueConfig.ProviderAlias]
		if !found {
			fmt.Printf("warning: no provider found for '%s'\n", valueConfig.ProviderAlias)
			continue
		}

		value, err := provider.GetValue(valueConfig.ValueUrlRaw)
		if err != nil {
			fmt.Printf("error: no key found for '%s' in provider '%s'\n", valueConfig.ValueUrl, valueConfig.ProviderAlias)
			os.Exit(1)
		}

		os.Setenv(valueConfig.EnvKey, value)

		values[valueConfig.EnvKey] = value
	}

	return values
}

func parseValueConfig(rawKeyValue string) (ValueConfig, bool) {
	rawPair := strings.SplitN(rawKeyValue, "=", 2)
	envKey, value := rawPair[0], rawPair[1]
	if !strings.HasPrefix(value, "nv://") {
		// other format checks, such as allowed chars, should be done
		return ValueConfig{}, false
	}

	valueUrl, err := url.Parse(value)
	if err != nil {
		fmt.Printf(
			"url parse error: %s: %s\n",
			value,
			err,
		)
		os.Exit(1)
		return ValueConfig{}, false
	}

	valueConfig := ValueConfig{
		ProviderAlias: valueUrl.Host,
		ValueUrlRaw:   value,
		ValueUrl:      valueUrl,
		EnvKey:        envKey,
	}

	return valueConfig, true
}
