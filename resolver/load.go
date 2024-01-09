package resolver

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v3"

	"github.com/maxnorth/nv/providers"
	"github.com/maxnorth/nv/providers/commandprovider"
	"github.com/maxnorth/nv/providers/commandprovider/commandgetprovider"
	"github.com/maxnorth/nv/providers/commandprovider/commandlistprovider"
	"github.com/maxnorth/nv/providers/hcvaultprovider"
	"github.com/maxnorth/nv/providers/sopsprovider"
)

func Load() *Resolver {
	r := &Resolver{
		providers: map[string]providers.Provider{},
	}

	yamlBytes, err := os.ReadFile("./nv.yaml")
	if err != nil {
		fmt.Println("nv.yaml not found")
		os.Exit(1)
	}
	jsonBytes := yamlToJson(yamlBytes)

	// needs validation
	gjson.GetBytes(jsonBytes, "resolvers").ForEach(func(key, value gjson.Result) bool {
		providerAlias := key.String()

		if value.Type == gjson.String {
			r.providers[providerAlias] = commandlistprovider.New(commandprovider.Config{
				Command: value.String(),
			})
			return true
		}

		providerType := "command"
		if providerProp := value.Get("provider"); providerProp.Exists() {
			providerType = providerProp.String()
		}

		switch providerType {
		case "command":
			var config commandprovider.Config
			if err := json.Unmarshal([]byte(value.Raw), &config); err != nil {
				panic(err)
			}
			switch config.Mode {
			case "":
				fallthrough
			case "list":
				r.providers[providerAlias] = commandlistprovider.New(config)
			case "get":
				r.providers[providerAlias] = commandgetprovider.New(config)
			default:
				fmt.Printf(
					"nv.yaml error: command resolver '%s' has unrecognized mode '%s'\n",
					providerAlias,
					providerType,
				)
				os.Exit(1)
			}
		case "hc-vault":
			var config hcvaultprovider.Config
			if err := json.Unmarshal([]byte(value.Raw), &config); err != nil {
				panic(err)
			}
			r.providers[providerAlias] = hcvaultprovider.New(config)
		case "sops":
			var config sopsprovider.Config
			if err := json.Unmarshal([]byte(value.Raw), &config); err != nil {
				panic(err)
			}
			r.providers[providerAlias] = sopsprovider.New(config)
		default:
			fmt.Printf(
				"nv.yaml error: resolver '%s' has unrecognized provider '%s'\n",
				providerAlias,
				providerType,
			)
			os.Exit(1)
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
