# Resolvers

Resolvers are used to convert references to real values, typically by fetching them from an external provider like Vault.

I want to try using this. How do I?
I'm going to install it and try running it, see what happens.
I might skim down the top part of the docs.
I might end up in the reference if I'm getting more serious about setting this up.

```yaml
resolvers:
  name:
    type: commmand
    command: echo "EXAMPLE_KEY=example-value"
```

### Resolver types

Each resolver definition has a required `type` option specifying which resolver to use. The rest of the options are dependent on the configured `type`.

- **type** string _required_ (**accepted values**: `command`, `hc-vault`, `sops`)
  <br/> The command to run. Value pairs will be extracted from stdout and parsed using the format specified in **output**. Will be evaluated using the default shell in `$SHELL`. Supports multi-line scripts.

#### command

##### Options:

- **run** string _required_
  <br/> The command to run. Value pairs will be extracted from stdout and parsed using the format specified in **output**. Will be evaluated using the default shell in `$SHELL`. Supports multi-line scripts.

- **output** string _optional_ (**accepted values**: `dotenv`, `json`, `yaml`; **default:** `dotenv`)
  <br/> The format of the output produced by the command.
