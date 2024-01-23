package resolver

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v3"

	"github.com/maxnorth/nv/providers"
	"github.com/maxnorth/nv/providers/commandprovider"
)

func Load() *Resolver {
	err := godotenv.Load(".env", ".env.local")
	if err != nil {
		panic(err)
	}

	r := &Resolver{
		providers:       map[string]providers.Provider{},
		loadedProviders: map[providers.Provider]struct{}{},
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
			r.providers[providerAlias] = commandprovider.New(commandprovider.Config{
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
			r.providers[providerAlias] = commandprovider.New(config)
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
