package commandprovider

import (
	"os/exec"
	"strings"
)

func (p *provider) GetValue(key string) (string, error) {
	subKeys := strings.Split(key, ":")

	args := []string{"-c", p.config.Command}
	args = append(args, subKeys...)

	output, err := exec.Command("bash", args...).Output()

	return string(output), err
}

func (p *provider) Load() error {
	if len(p.config.Command) == 0 {
		panic("command is not defined")
	}

	return nil
}
