package sopsprovider

import (
	"encoding/json"
	"os/exec"

	"gopkg.in/yaml.v3"
)

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
