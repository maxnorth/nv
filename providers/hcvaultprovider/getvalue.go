package hcvaultprovider

import (
	"errors"
)

func (p *provider) GetValue(key string) (string, error) {
	value, exists := p.values[key]
	if !exists {
		return "", errors.New("key not found in vault")
	}

	return value, nil
}
