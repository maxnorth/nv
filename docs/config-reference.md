# Resolvers

Resolvers are used to convert value references to real values, typically by fetching them from an external provider like Vault.

I want to try using this. How do I?
I'm going to install it and try running it, see what happens.
I might skim down the top part of the docs.
I might end up in the reference if I'm getting more serious about setting this up.

```yaml
resolvers:
  [provider-alias]:
    kind: commmand
    command: string
    profiles:
      default:
        command: string
      production:
        command: string
      staging:
        command:
```

```yaml
command: [string] The shell command to run. Evaluated using `sh`.
mode: [string] get|list. Defaults to 'list'.
output: [string] dotenv|yaml|json. Defaults to 'dotenv'.
```

**command** string _required_
<br/> Here's a description of this thing. Here's a description of this thing. Here's a description of this thing. Here's a description of this thing. Here's a description of this thing. Here's a description of this thing.

**mode** string (get|list) _optional_ (default: list)
<br/> Here's a description of this thing. Here's a description of this thing. Here's a description of this thing. Here's a description of this thing. Here's a description of this thing. Here's a description of this thing.

**output** string (**enum**: dotenv, json, yaml) _optional_ (**default**: dotenv)
<br/> Here's a description of this thing. Here's a description of this thing. Here's a description of this thing. Here's a description of this thing. Here's a description of this thing.
