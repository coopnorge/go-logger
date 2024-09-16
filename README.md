# go-logger

This Go package is used to offer a unified logging interface among projects.

## Documentation

There usage documentation in [docs](docs/). Development documentation is in
this file.

## Development workflow

### Validate

```bash
docker compose run --rm golang-devtools validate
```

### Other targets

```bash
docker compose run --rm golang-devtools help
```

## Mocks

To generate or update mocks use
[`gomockhandler`](https://github.com/sanposhiho/gomockhandler). `gomockhandler`
is provided by `golang-devtools`.

### Check mocks

```bash
docker compose run --rm golang-devtools gomockhandler -config ./gomockhandler.json check
```

### Generate / Update mocks

```bash
docker compose run --rm golang-devtools gomockhandler -config ./gomockhandler.json mockgen
```

## User documentation

User documentation is build using TechDocs and published to
[Inventory](https://inventory.internal.coop/docs/default/component/go-logger).

To list the commands available for the TechDocs image:

```sh
docker compose run --rm help
```

For more information see the [TechDocs Engineering
Image](https://github.com/coopnorge/engineering-docker-images/tree/main/images/techdocs).

### Documentation validation

To Validate changed documentation:

```sh
docker compose run --rm techdocs validate
```

To validate all documentation:

```sh
docker compose run --rm techdocs validate MARKDOWN_FILES=docs/
```

### Documentation preview

To preview the documentation:

```sh
docker compose up techdocs
```
