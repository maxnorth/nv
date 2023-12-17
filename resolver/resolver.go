package resolver

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/maxnorth/nv/providers"
)

type ValueConfig struct {
	ProviderAlias string
	SourceKey     string
	TargetKey     string
}

func Resolve(providers map[string]providers.Provider) {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// load providers eagerly - start dumb
	// check environment for providers and keys
	valueConfigs := []ValueConfig{}
	for _, e := range os.Environ() {
		if valueConfig, found := parseValueConfig(e); found {
			valueConfigs = append(valueConfigs, valueConfig)
		}
	}

	for _, valueConfig := range valueConfigs {
		provider, found := providers[valueConfig.ProviderAlias]
		if !found {
			fmt.Printf("no provider found for '%s'\n", valueConfig.ProviderAlias)
			continue
		}

		value, err := provider.GetValue(valueConfig.SourceKey)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%s=\"%s\"\n", valueConfig.TargetKey, value)
	}
}

func parseValueConfig(rawKeyValue string) (ValueConfig, bool) {
	rawPair := strings.SplitN(rawKeyValue, "=", 2)
	targetKey, rawSource := rawPair[0], rawPair[1]
	if !strings.HasPrefix(rawSource, "@") {
		// other format checks, such as allowed chars, should be done
		return ValueConfig{}, false
	}

	trimmedSource := strings.TrimPrefix(rawSource, "@")
	sourcePair := strings.SplitN(trimmedSource, ":", 2)

	providerAlias := sourcePair[0]
	var sourceKey string
	if len(sourcePair) == 2 {
		sourceKey = sourcePair[1]
	} else {
		sourceKey = targetKey
	}

	valueConfig := ValueConfig{
		ProviderAlias: providerAlias,
		SourceKey:     sourceKey,
		TargetKey:     targetKey,
	}

	return valueConfig, true
}
