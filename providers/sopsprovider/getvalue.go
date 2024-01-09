package sopsprovider

import (
	"errors"

	"github.com/tidwall/gjson"
)

func (p *provider) GetValue(key string) (string, error) {
	jsonResult := gjson.GetBytes(p.jsonData, key)
	exists := jsonResult.Exists()
	if !exists {
		return "", errors.New("key not found in sops file")
	}

	return jsonResult.String(), nil
}
