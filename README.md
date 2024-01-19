# nv - Environment variables with reach

**nv** is an environment variable loader and URL resolver. It was created to improve the developer experience of working with secret providers.

## How it works

**nv** allows environment variables to contain special URL's that locate values and automatically resolves them at runtime. It also runs the `dotenv` loader so you can keep values in `.env` files.

```dotenv
# .env
DB_USER=nv://vault/kv/data/my-service?field=db-user
DB_PASSWORD=nv://vault/kv/data/my-service?field=db-password
VENDOR_API_KEY=nv://vault/kv/data/my-service?field=vendor-key

# .env.local
DB_USER=postgres
DB_PASSWORD=postgres
VENDOR_API_KEY=nv://1password/development/my-vendor/keys/api-key
```

Several providers are supported out of the box. Connections use common defaults, or can be configured using an `nv.yaml` file.

```yaml
# nv.yaml
resolvers:
  vault:
    type: hc-vault
    host: http://localhost:8200
  onepassword:
    type: 1password
```

**nv** is designed to be usable purely as a CLI, but it is also available as a golang library. Libraries for more languages may be added depending on interest.

## Why use nv

- Decouples applications from config/secret providers, and saves you from writing code to use them
- Increases project consistency by enabling one solution for every environment, workflow, and application
- Makes using shared secrets in local development easy and secure by default
- Results in centralized and self-documenting use of external config providers

# Getting started

### Installation

<!-- need a solution for distributing the CLI -->

### Basic usage

`nv` is as an application launcher. It will execute the provided command with variables loaded into the environment:

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

1. Load: read environment variables from local .env files, like `dotenv` does
2. Resolve: check all environment variables for value references, for example: `@vault:kv/my-app/database/password`. `nv` resolves the real value using the provider config in `nv.yaml` and rewrites the environment variable.

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
