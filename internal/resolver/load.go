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
	err := runDotenv(env)
	if err != nil {
		return nil, err
	}

	r := &Resolver{
		providers:       map[string]providers.Provider{},
		loadedProviders: map[providers.Provider]struct{}{},
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
				"nv.yaml error: resolver '%s' has unrecognized provider '%s'\n",
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

func runDotenv(env string) error {
	files := []string{".env", fmt.Sprintf(".env.%s", env)}

	for _, f := range files {
		switch err := godotenv.Load(f).(type) {
		case nil:
		case *fs.PathError:
		default:
			return fmt.Errorf("error: failed to load %s file: %s\n", f, err)
		}
	}

	return nil
}
