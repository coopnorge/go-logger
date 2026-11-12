# Releasing

The repository holds nine Go modules that are released together with the same
version, as in [opentelemetry-go](https://github.com/open-telemetry/opentelemetry-go/blob/main/RELEASING.md):

| Module                                              | Directory           | Tag                        |
| --------------------------------------------------- | ------------------- | -------------------------- |
| `github.com/coopnorge/go-logger`                    | `.`                 | `vX.Y.Z`                   |
| `github.com/coopnorge/go-logger/adapter/echo.v4`    | `adapter/echo.v4`   | `adapter/echo.v4/vX.Y.Z`   |
| `github.com/coopnorge/go-logger/adapter/goose.v3`   | `adapter/goose.v3`  | `adapter/goose.v3/vX.Y.Z`  |
| `github.com/coopnorge/go-logger/adapter/gorm`       | `adapter/gorm`      | `adapter/gorm/vX.Y.Z`      |
| `github.com/coopnorge/go-logger/adapter/kratos.v2`  | `adapter/kratos.v2` | `adapter/kratos.v2/vX.Y.Z` |
| `github.com/coopnorge/go-logger/adapter/echo`       | `adapter/echo`      | `adapter/echo/vX.Y.Z`      |
| `github.com/coopnorge/go-logger/adapter/goose`      | `adapter/goose`     | `adapter/goose/vX.Y.Z`     |
| `github.com/coopnorge/go-logger/adapter/kratos`     | `adapter/kratos`    | `adapter/kratos/vX.Y.Z`    |

The last three are deprecated forwarder modules, see below.

`e2e_tests` is a separate module that is never released.

## One module per library major version

An adapter for a library whose import path has a major version above 1 gets its
own module for each major version of that library. The major version is a dotted
suffix of the directory name: `adapter/echo.v4` for `github.com/labstack/echo/v4`,
`adapter/goose.v3` and `adapter/kratos.v2`. The module path is the directory
path, so the tags look like `adapter/echo.v4/vX.Y.Z`. This is the same layout
as [dd-trace-go](https://github.com/DataDog/dd-trace-go), for example
`contrib/labstack/echo.v4`.

Do not use a `/v4` path element. A module path that ends in `/vN` requires
module version vN
([Go modules reference](https://go.dev/ref/mod#major-version-suffixes)), which
clashes with our v0.x tags. `.v4` is not a `/vN` suffix, so the module keeps the
shared v0.x version. A later Echo major becomes a sibling module, `adapter/echo.v5`,
so each major has its own dependencies and can be deprecated separately with a
`// Deprecated:` line in its `go.mod`. If go-logger itself gets a v2, all
modules get a `/v2` suffix, as dd-trace-go did.

The package name does not have to match the last path element. `echo.v4` uses
`package echo`, `goose.v3` uses `package goose` and `kratos.v2` uses
`package logs`. Importers should set an alias, for example
`echov4 "github.com/coopnorge/go-logger/adapter/echo.v4"`.

Adapters for libraries without a major version suffix, such as Gorm, stay at
`adapter/<name>`.

## Deprecated forwarder modules

`adapter/echo`, `adapter/goose` and `adapter/kratos` are small forwarder
modules that `require` the versioned module. Each `go.mod` has a
`// Deprecated: use github.com/coopnorge/go-logger/adapter/<name>.vN` line above
`module`, so `go list -m -u` and editors warn consumers. Each exported identifier
is a type alias or a thin wrapper annotated with `//go:fix inline`, so
`go fix -inline ./...` migrates users to the versioned package. Remove the
forwarders at the next breaking release.

`versions.yaml` lists the modules in the `go-logger` module set.
[multimod](https://github.com/open-telemetry/opentelemetry-go-build-tools/tree/main/multimod)
reads it, so adding an adapter means adding its module path to that list.

Each adapter's `go.mod` requires the core module and has a committed
`replace github.com/coopnorge/go-logger => ../..` directive, so development and
CI always build against the core in the same commit. Consumers ignore the
`replace` directive and get the version in the `require` line. That line is
updated by `mage prerelease`; Renovate is configured not to touch it.

## Manual release

1. Pick the next version and edit `version` in `versions.yaml`. While the major
   version is 0, use a minor bump for breaking changes and a patch bump
   otherwise.
2. Make sure you are on an up-to-date, clean `main`, then run:

   ```bash
   go tool mage prerelease
   ```

   This creates the branch `prerelease_go-logger_vX.Y.Z`, sets the version of
   the core module in every adapter's `go.mod`, runs `go mod tidy` and commits.
3. Push the branch and open a pull request. CI must pass.
4. Merge the pull request. Use a merge or squash commit, then pull `main` and
   note the hash of the merge commit.
5. Tag all modules on that commit:

   ```bash
   go tool mage tag "$(git rev-parse HEAD)"
   ```

   Run it only when `HEAD` is the merge commit. If tagging fails halfway,
   multimod removes the tags it created.
6. Push all the tags, including the ones for the adapters:

   ```bash
   git push origin vX.Y.Z \
     adapter/echo.v4/vX.Y.Z adapter/goose.v3/vX.Y.Z adapter/gorm/vX.Y.Z adapter/kratos.v2/vX.Y.Z \
     adapter/echo/vX.Y.Z adapter/goose/vX.Y.Z adapter/kratos/vX.Y.Z
   ```

7. Publish the draft release created by release-drafter, using the tag `vX.Y.Z`.
8. Check that the Go proxy serves every module, for example
   `go list -m github.com/coopnorge/go-logger/adapter/gorm@vX.Y.Z`.

A wrongly pushed tag cannot be moved, because the Go proxy caches it. Release a
new patch version and add a `retract` directive for the bad version to the
`go.mod` of the affected module, see the deprecation section.

## Consistency check

```bash
go tool mage consistency
```

This runs `crosslink` to sync the `replace` directives, `multimod verify` to
check `versions.yaml`, and `go mod tidy` in every module. It fails if that
changes any `go.mod`, `go.sum` or `versions.yaml` file. Run it after changing
dependencies.

## Deprecating and removing an adapter

1. Add a `// Deprecated: <reason and replacement>` comment on the line before
   `module` in the adapter's `go.mod`. Release as above.
2. Keep the adapter for at least one more tagged release so consumers see the
   notice.
3. Remove the directory and its entry in `versions.yaml`. Old tags stay
   available in the Go proxy.

To mark a version as broken, add `retract vX.Y.Z // reason` to the `go.mod` of
the affected module and release a new version.

## Future: rolling release

The release is manual today. The setup is prepared so that a workflow can run
the same steps without a person:

- Run every two weeks and on demand.
- Bump `versions.yaml`, run `mage prerelease`, open a release pull request and
  auto-merge it when CI passes.
- After the merge, a bot runs `mage tag` on the merge commit, pushes all tags
  and publishes the release notes.

Prerequisites:

- A policy exception that allows the bot to merge its own pull request.
- A token that can push tags. The default `GITHUB_TOKEN` does not trigger other
  workflows, and tag rules may block it.
