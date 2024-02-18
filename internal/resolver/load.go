package resolver

import (
	"encoding/json"
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
		configFound:     true,
		LoadedVars:      keys,
	}

	yamlBytes, err := os.ReadFile("./nv.yaml")
	if _, isPathErr := err.(*fs.PathError); err != nil && !isPathErr {
		return nil, fmt.Errorf("failed to load nv.yaml: %s", err)
	} else if err != nil {
		r.configFound = false
	}

	jsonBytes, err := yamlToJson(yamlBytes)
	if err != nil {
		return nil, err
	}

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

func yamlToJson(yamlBytes []byte) ([]byte, error) {
	var result any
	err := yaml.Unmarshal(yamlBytes, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse nv.yaml: %s", err)
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}

	return jsonBytes, nil
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
