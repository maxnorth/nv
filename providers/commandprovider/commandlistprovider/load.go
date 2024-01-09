package commandlistprovider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

func (p *provider) Load() error {
	if len(p.config.Command) == 0 {
		fmt.Println("command is not defined")
		os.Exit(1)
	}

	output, err := exec.Command("sh", "-c", p.config.Command).Output()
	if err != nil {
		panic(err)
	}

	switch p.config.Output {
	case "":
		fallthrough
	case "dotenv":
		return p.parseDotenv(output)
	case "json":
		return p.parseJson(output)
	case "yaml":
		return p.parseYaml(output)
	}

	return nil
}

func (p *provider) parseDotenv(output []byte) error {
	reader := bytes.NewReader(output)
	values, err := godotenv.Parse(reader)
	if err != nil {
		panic(err)
	}

	p.values = values

	return nil
}

func (p *provider) parseJson(output []byte) error {
	err := json.Unmarshal(output, &p.values)
	if err != nil {
		return err
	}

	return nil
}

func (p *provider) parseYaml(output []byte) error {
	err := yaml.Unmarshal(output, &p.values)
	if err != nil {
		return err
	}

	return nil
}
