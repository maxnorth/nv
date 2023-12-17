package commandprovider

import "errors"

func (cp *provider) GetValue(key string) (string, error) {
	value, exists := cp.values[key]
	if !exists {
		return "", errors.New("key not found in command output")
	}

	return value, nil
}
