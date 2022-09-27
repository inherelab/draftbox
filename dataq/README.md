# Data query

`dq` like `jq` command, but can use for YAML, JSON, TOML, INI, .env

## Usage


```shell
dq '{json string}' | dq .top
```

Read input from file:

```shell
dq path.json | dq .top
```

Read input from Clipboard:

```shell
dq @c | dq .top
```

