package resolver

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v3"

	"github.com/maxnorth/nv/providers"
	"github.com/maxnorth/nv/providers/commandprovider"
	"github.com/maxnorth/nv/providers/hcpvaultsecretsprovider"
	"github.com/maxnorth/nv/providers/sopsprovider"
)

func Load() *Resolver {
	r := &Resolver{
		providers: map[string]providers.Provider{},
	}

	yamlBytes, err := os.ReadFile("./nv.yaml")
	if err != nil {
		fmt.Println("nv.yaml not found")
		panic(err)
	}
	jsonBytes := yamlToJson(yamlBytes)

	// needs validation
	gjson.GetBytes(jsonBytes, "renderers").ForEach(func(key, value gjson.Result) bool {
		providerAlias := key.String()

		if value.Type == gjson.String {
			r.providers[providerAlias] = commandprovider.New(commandprovider.Config{
				Command: value.String(),
			})
			return true
		}

		providerType := value.Get("type").String()
		switch providerType {
		case "command":
			var config commandprovider.Config
			if err := json.Unmarshal([]byte(value.Raw), &config); err != nil {
				panic(err)
			}
			r.providers[providerAlias] = commandprovider.New(config)
		case "hashicorp-vault-secrets":
			var config hcpvaultsecretsprovider.Config
			if err := json.Unmarshal([]byte(value.Raw), &config); err != nil {
				panic(err)
			}
			r.providers[providerAlias] = hcpvaultsecretsprovider.New(config)
		case "sops":
			var config sopsprovider.Config
			if err := json.Unmarshal([]byte(value.Raw), &config); err != nil {
				panic(err)
			}
			r.providers[providerAlias] = sopsprovider.New(config)
		}

		return true
	})

	return r
}

func yamlToJson(yamlBytes []byte) []byte {
	var result any
	yaml.Unmarshal(yamlBytes, &result)

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}

	return jsonBytes
}
