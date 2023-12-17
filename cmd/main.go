package main

import (
	"github.com/maxnorth/nv/providers"
	commandprovider "github.com/maxnorth/nv/providers/command"
	"github.com/maxnorth/nv/resolver"
)

// func main() {
// 	app := &cli.App{
// 		Name:      "boom",
// 		Usage:     "make an explosive entrance",
// 		UsageText: "whoooop",
// 		Action: func(c *cli.Context) error {

// 			return nil
// 		},
// 		Commands: []*cli.Command{
// 			{},
// 		},
// 	}

// 	if err := app.Run(os.Args); err != nil {
// 		log.Fatal(err)
// 	}
// }

func main() {

	type providerDef struct {
		Type   string
		Config any
	}

	providerDefs := map[string]providerDef{
		"customfile": {
			Type: "command",
			Config: commandprovider.Config{
				Command: "cat",
				Args:    []string{"customfile"},
			},
		},
		"vault": {
			Type: "hashicorp-vault",
			Config: commandprovider.Config{
				Command: "cat",
				Args:    []string{"vaultfile"},
			},
		},
	}

	providers := map[string]providers.Provider{}
	for providerAlias, providerDef := range providerDefs {
		switch providerDef.Type {
		case "command":
			providers[providerAlias] = commandprovider.New(
				providerDef.Config.(commandprovider.Config),
			)
		case "hashicorp-vault":
			// TODO: add actual vault provider
			providers[providerAlias] = commandprovider.New(
				providerDef.Config.(commandprovider.Config),
			)
		}
	}

	for _, provider := range providers {
		provider.Load()
	}

	resolver.Resolve(providers)
}

// replace env vars starting with @
//  - resolve the provider
//    - if provider not found, fail w/ helpful message
//    - later provide an option to ignore specific values or any not found value
//  - if no : the key is the map target
//    - if : use that as the key
//  - using the key, invoke the provider and resolve the value
//  - map the value to the target env var

// need to identify resolution order and ordering issues
