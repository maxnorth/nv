package resolver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/joho/godotenv"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v3"

	"github.com/maxnorth/nv/internal/providers"
	"github.com/maxnorth/nv/internal/providers/commandprovider"
)

func Load(env string) (*Resolver, error) {
	keys, err := runDotenv(env)
	if err != nil {
		return nil, err
	}

	r := &Resolver{
		providers:       map[string]providers.Provider{},
		loadedProviders: map[providers.Provider]struct{}{},
		LoadedVars:      keys,
	}

	yamlBytes, err := os.ReadFile("./nv.yaml")
	if err != nil {
		return nil, errors.New("nv.yaml not found")
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
			err = fmt.Errorf(
				"nv.yaml resolver '%s' has unrecognized provider '%s'\n",
				providerAlias,
				providerType,
			)
			return false
		}

		return true
	})

	if err != nil {
		return nil, err
	}

	return r, nil
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

func runDotenv(env string) ([]string, error) {
	files := []string{fmt.Sprintf(".env.%s", env), ".env"}

	keys := []string{}
	for _, f := range files {
		envMap, err := godotenv.Read(f)
		switch err := err.(type) {
		case nil:
		case *fs.PathError:
		default:
			return nil, fmt.Errorf("failed to load %s file: %s", f, err)
		}

		for key, val := range envMap {
			if _, exists := os.LookupEnv(key); !exists {
				keys = append(keys, key)
				os.Setenv(key, val)
			}
		}
	}

	return keys, nil
}
