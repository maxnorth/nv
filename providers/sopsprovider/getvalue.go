package sopsprovider

import (
	"encoding/json"
	"errors"
	"os/exec"

	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v3"
)

func (p *provider) GetValue(key string) (string, error) {
	jsonResult := gjson.GetBytes(p.jsonData, key)
	exists := jsonResult.Exists()
	if !exists {
		return "", errors.New("key not found in sops file")
	}

	return jsonResult.String(), nil
}

func (p *provider) Load() error {
	if len(p.config.File) == 0 {
		panic("sops file is not defined")
	}

	output, err := exec.Command("sops", "-d", p.config.File).Output()
	if err != nil {
		// TODO: wrap
		return err
	}

	var values any
	// TODO multiple possible formats?
	err = yaml.Unmarshal(output, &values)
	if err != nil {
		// TODO: wrap
		return err
	}

	jsonData, err := json.Marshal(values)
	if err != nil {
		// TODO: wrap
		return err
	}

	p.jsonData = jsonData

	return nil
}
