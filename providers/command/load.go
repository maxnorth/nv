package commandprovider

import (
	"bytes"
	"os/exec"

	"github.com/joho/godotenv"
)

func (p *provider) Load() error {
	// run the command and capture the output values
	output, err := exec.Command(p.config.Command, p.config.Args...).Output()
	if err != nil {
		// TODO: wrap
		return err
	}

	p.values, err = godotenv.Parse(bytes.NewReader(output))
	if err != nil {
		// TODO: wrap
		return err
	}

	return nil
}
