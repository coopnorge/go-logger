# go-logger

This Go package is used to offer a unified logging interface among projects.

## Documentation

There usage documentation in [docs](docs/). Development documentation is in this
file.

## Development workflow

The source code is build using `mage`.

### Repository layout

The core module is in the repository root. Each adapter in `adapter/<name>` is
a separate Go module that is released together with the core module, see
[RELEASING.md](RELEASING.md). Each adapter's `go.mod` has a `replace` directive
that points to the core module in the same checkout.

Adapters for libraries with a major version above 1 get one module per library
major version, named with a dotted suffix: `adapter/echo.v4`, `adapter/goose.v3`
and `adapter/kratos.v2`. `adapter/echo`, `adapter/goose` and `adapter/kratos` are
deprecated forwarder modules. See [RELEASING.md](RELEASING.md).

`mage go:test`, `mage go:lint` and the other `go:` targets run for every
`go.mod` found in the repository.

### Prerequisites

1. Install Go version 1.24 or later and
   [Docker](https://docs.docker.com/get-docker/).

2. Install Go tools:

   ```bash
   go install tool
   ```

### Validate

```bash
go tool mage validate
```

### Check modules are consistent

After changing dependencies, run:

```bash
go tool mage consistency
```

### Other targets

```bash
go tool mage -l
```

## Releasing

See [RELEASING.md](RELEASING.md).

## User documentation

User documentation is build using TechDocs and published to
[Inventory](https://inventory.internal.coop/docs/default/component/go-logger).

To list the commands available for the TechDocs image:

```sh
docker compose run --rm help
```

For more information see the
[TechDocs Engineering Image](https://github.com/coopnorge/engineering-docker-images/tree/main/images/techdocs).

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
