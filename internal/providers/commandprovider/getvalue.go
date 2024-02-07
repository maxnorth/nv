package commandprovider

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

// TODO: standardize on URL's in interface
func (p *provider) GetValue(rawUrl string) (string, error) {
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}

	cmd := exec.Command("sh", "-c", p.config.Command)

	cmd.Env = append(os.Environ(), getUrlVars(parsedUrl)...)

	var stderr bytes.Buffer
	// cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	output, err := cmd.Output()

	if err != nil {
		return "", fmt.Errorf("error in '%s' resolver for url '%s':\n%s\n", parsedUrl.Host, rawUrl, stderr.Bytes())
	}

	trimmedOutput := strings.Trim(string(output), "\n")

	return trimmedOutput, err
}

func getUrlVars(refUrl *url.URL) []string {
	result := []string{}

	path := refUrl.Host + refUrl.Path
	result = append(result, "NV_URL="+refUrl.String(), "NV_URL_PATH="+path)

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
			result = append(result, "NV_URL_ARG_"+key+"="+value)
		}
	}

	return result
}
