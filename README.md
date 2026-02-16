# go-logger

This Go package is used to offer a unified logging interface among projects.

## Documentation

There usage documentation in [docs](docs/). Development documentation is in this
file.

## Development workflow

The source code is build using `mage`.

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

### Other targets

```bash
go tool mage -l
```

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
