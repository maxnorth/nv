package commandprovider

import (
	"fmt"
	"os"
)

func (p *provider) Load() error {
	if len(p.config.Command) == 0 {
		fmt.Println("command is not defined")
		os.Exit(1)
	}

	return nil
}
