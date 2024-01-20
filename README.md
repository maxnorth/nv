# nv - Enviably simple app configuration

**nv** is a dotenv-style environment variable loader that makes secret managers easier to work with. No more copying secrets to developer machines or bending [12FA principles](https://12factor.net/config) by tightly coupling services to secret managers.

## Intro

**nv** looks for environment variables containing `nv://` URL's and automatically resolves them to a value. It also runs the `dotenv` loader so you can manage values by environment in `.env` files.

```dotenv
# .env
DB_PASSWORD=nv://vault/kv/data/my-service?field=db-password
VENDOR_API_KEY=nv://vault/kv/data/my-service?field=vendor-key
```

```
# .env.local - easily override with static values or other managers
DB_PASSWORD=postgres
VENDOR_API_KEY=nv://1password/development/my-vendor/keys/api-key
```

These URLs are matched to **resolvers**, commands that **nv** runs to obtain a value for the URL, which you define in `nv.yaml`.

```yaml
# NV_URL_* env vars are provided to access parts of the URL
resolvers:
  vault: vault kv get -mount="secret" -field=$NV_URL_ARG_FIELD $NV_URL_PATH
  sops: sops -d --extract $NV_URL_ARG_EXTRACT $NV_URL_PATH
  1password: op read op://$NV_URL_PATH
```

The primary `nv` commands are **run** and **print**.

```bash
$ nv run -- npm start           # run a command with env vars loaded and resolved.
$ nv -- npm start               # 'run' is the default and can be skipped for convenience.
$ nv -- zsh                     # it can be handy to start a shell session with vars loaded.
$ nv -e staging -- zsh          # run using a targeted environment.

$ nv print                      # load, resolve, and print values for inspection.
DB_PASSWORD=hello-from-vault!
VENDOR_API_KEY=hello-again!
$ nv print --output json        # print in json or yaml for easy parsing by applications.
{
  "DB_PASSWORD": "hello-from-vault!",
  "VENDOR_API_KEY": "hello-again!"
}
```

## Installation

<!-- need a solution for distributing the CLI -->

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
