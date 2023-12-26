# nv

`nv` is an environment variable loader, similar to `dotenv` but with a few differences:

- `nv` is not limited to .env files, it can also load env vars using config and secret management providers like Vault, `sops`, and more
- nv is a portable CLI tool, not a language dependent library

# Goals

The goals of this project are to:

- Provide a well-encapsulated abstraction for env variable sourcing that is decoupled from application code
- Make using secrets in local development easy and secure by default
- Be a single tool suited to any tech stack and execution environment

# How it works

nv is used to run your applications with environment variables loaded from anywhere you need. It does so by resolving directives stored in environment variables, which are instructions on where to fetch the required value.

```
$ export MY_SECRET=@vault:kv/secrets/my-secret
$ nv -- printenv MY_SECRET
Hello from vault!
```

In this example, `vault` is an alias for a provider defined in a `nv.yaml` file in the current directory.

```
# nv.yaml
providers:
    vault:
        type: hcl-vault
        host: http://localhost:8200
```

nv comes with built-in support a few common providers. For unsupported cases, you can extend it using the `command` provider:

# Installation

- install it

# Quickstart
