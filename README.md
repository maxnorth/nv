# nv - Environment variables with reach

**nv** is an environment variable loader and URL resolver. It aims to fulfill an unmet promise of secret managers, which should make it **easy** to use application secrets **securely** in all environments, especially in local development.

- How it works
- Installation
- Command usage
- Resolver config
- FAQ

## How it works

**nv** looks for environment variables containing special URL's that locate values and automatically resolves them at runtime. It also runs the `dotenv` loader so you can keep values in `.env` files and easily manage differences in environments.

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

Resolvers allow `nv` to obtain values for a URL. You define them in a `nv.yaml` file typically kept at the root of the project.

```yaml
resolvers:
  vault: vault kv get -mount="secret" -field=$NV_URL_ARG_FIELD $NV_URL_PATH
  sops: sops -d --extract $NV_URL_ARG_EXTRACT $NV_URL_PATH
  1password: op read op://$NV_URL_PATH
```

**nv** is designed to be usable purely as a CLI. There are two essential `nv` commands, **run** and **print**.

```bash
$ nv run -- node dist/main.js   # run a command with env vars loaded and resolved.
$ nv -- node dist/main.js       # 'run' is the default and can be dropped if '--' is present.
$ nv -- zsh                     # conveniently start a shell session with vars loaded.
$ exit                          # exit when done to reset variables.
$ nv -p staging -- zsh          # then load values for a different environment instead.
```

```bash
$ nv print      # resolve and print all values loaded from .env files

$ nv print -o json

$ nv -- zsh                   # a convenient local dev pattern is to start a shell with vars loaded

$ nv -p staging -- zsh        # or start a shell with values for a different environment instead
```

The CLI approach allows `nv` to be used in any runtime environment and any language or tech stack. This also makes it usable with software you don't own the code for.

Language-specific libraries are not currently available, but could bring modest convenience, and may be explored depending on interest.

## Installation

<!-- need a solution for distributing the CLI -->

## Command usage

```bash
$ nv -- node dist/main.js   # launch an application with env vars loaded
```

OR

```bash
$ nv -- zsh                 # launch a new shell session with env vars loaded
$ node dist/main.js
```

## Resolver config

Each resolver has a name and a command to run to resolve URLs. Commands are called once per URL, and the output of the command is used as the value. URLs are matched to resolvers by the host portion of the URL (see table below if URL portion clarity is needed).

```yaml
resolvers:
  vault: vault kv get -mount="secret" -field=$NV_URL_ARG_FIELD $NV_URL_PATH
  sops: sops -d --extract $NV_URL_ARG_EXTRACT $NV_URL_PATH
  1password: op read op://$NV_URL_PATH
```

As you can see above, special environment variables are made available to access parts of the URL. These will change on each run to match the current URL being resolved.

Here are all the variables provided for an example URL `nv://vault/something/or/other?field=thing`:

| Variable            | Value                                       |
| ------------------- | ------------------------------------------- |
| `$NV_URL`           | `nv://vault/something/or/other?field=thing` |
| `$NV_URL_HOST`      | `vault`                                     |
| `$NV_URL_PATH`      | `something/or/other`                        |
| `$NV_URL_ARG_FIELD` | `thing`                                     |

`$NV_URL_ARG_*` is dynamic and one will be provided for each query arg in the URL. Arg names will be uppercased in the variable, and only the first occurrence of an arg name/value pair will be used when duplicates/conflicts occur.

### Built-in resolvers

At this time, the custom command resolver described above is the only available option, and in theory it should be the only one needed.

That said, built-in resolvers for third party secret providers may be added to improve usability. Users are encouraged to open issues to request built-in support for a third party provider, especially if **A)** it's popularly used and **B)** the tooling needed has challenging installation requirements or usability issues.

### Adjusting resolver behavior by environment

There are a couple ways to change your resolver behavior by environment.

1. You can put variables in relevant `.env` files (`.env.production`, `.env.staging`) with different values, and use that variable in your resolver command. All `.env` files will be loaded before resolvers run.
2. If more control is needed than above, create a separate resolver entry for the different resolver versions and target the environment-appropriate resolver in the host portion of URLs in the relevant `.env` files.

When/if configuration for resolvers becomes more complex, additional ways to manage environment overrides will likely be introduced.
