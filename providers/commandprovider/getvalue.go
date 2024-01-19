package commandprovider

import (
	"net/url"
	"os"
	"os/exec"
	"strings"
)

// TODO: standardize on URL's in interface
func (p *provider) GetValue(key string) (string, error) {
	refUrl, err := url.Parse(key)
	if err != nil {
		return "", err
	}

	cmd := exec.Command("sh", "-c", p.config.Command)
	cmd.Env = os.Environ()
	path := strings.TrimLeft(refUrl.Path, "/")
	cmd.Env = append(cmd.Env, "NV_URL="+key, "NV_URL_HOST="+refUrl.Host, "NV_URL_PATH="+path)

	queryArgKeys := map[string]struct{}{}
	for key, values := range refUrl.Query() {
		// TODO: double check query chars and env chars compatible
		key = strings.ToUpper(key)
		value := ""
		if len(values) > 0 {
			value = values[0]
		}

		if _, exists := queryArgKeys[key]; !exists {
			queryArgKeys[key] = struct{}{}
			cmd.Env = append(cmd.Env, "NV_URL_ARG_"+key+"="+value)
		}
	}

	output, err := cmd.Output()

	trimmedOutput := strings.Trim(string(output), "\n")

	return trimmedOutput, err
}
