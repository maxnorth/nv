package commandlistprovider

import "errors"

func (p *provider) GetValue(key string) (string, error) {
	value, exists := p.values[key]
	if !exists {
		return "", errors.New("value not found")
	}

	return value, nil
}
