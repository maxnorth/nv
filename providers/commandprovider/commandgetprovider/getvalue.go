package commandgetprovider

import (
	"os/exec"
	"strings"
)

func (p *provider) GetValue(key string) (string, error) {
	subKeys := strings.Split(key, ":")

	args := []string{"-c", p.config.Command}
	args = append(args, subKeys...)

	output, err := exec.Command("sh", args...).Output()

	trimmedOutput := strings.Trim(string(output), "\n")

	return trimmedOutput, err
}
