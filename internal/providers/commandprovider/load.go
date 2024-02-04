package commandprovider

import "errors"

func (p *provider) Load() error {
	if len(p.config.Command) == 0 {
		return errors.New("command is not defined")
	}

	return nil
}
