# nv - simplified access to external configuration

The **nv** CLI is a unified interface for loading environment variables from anywhere.

It allows you to initialize environment variables with references to values stored in external providers like Vault, and then resolve them to their real values before your app starts.

<!-- GIF of usage here -->

## Why use **nv**

- Decouples applications from config/secret providers, and saves you time writing code to use them
- Reduces project complexity by providing a single pattern for every environment, workflow, and application
- Makes using shared secrets in local development easy and secure by default

# Getting started

### Installation

### Basic usage

The primary way to use `nv` is as an application launcher. It will execute the provided command with variables loaded into the environment:

```bash
$ nv -- node dist/main.js   # launch an application with env vars loaded
```

OR

```bash
$ nv -- zsh                 # launch a new shell session with env vars loaded
$ node dist/main.js
```

### How environment loading works

Environment loading consists of two steps:

1. Load: read environment variables from local .env files, same as `dotenv`
2. Resolve: check all environment variables for values matching a unique syntax. This syntax specifies a provider and location for a real value - for example: `@vault:kv/my-app/database/password`. `nv` resolves the real value using the provider, then assigns it back to the environment variable.

```
# .env
DB_PASSWORD=@vault:kv/my-app/database/password
```

```
$ nv -- printenv | grep DB_PASSWORD
DB_PASSWORD=This is my DB password stored in vault!
```

### How providers work

`nv` and its providers are configured in a `nv.yaml` file. `nv` will look for this file in the current working directory by default, and will fail if not found.

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

# FAQ

- Does my whole team need to adopt nv for me to use it?
- What if I want to load env vars in a way not supported by a built-in provider?
- Does it work for every OS?
