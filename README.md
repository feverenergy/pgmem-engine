# pgmem-engine

The compressed PostgreSQL binaries that [pgmem](https://github.com/feverenergy/pgmem) runs. pgmem is a Go library that starts a real PostgreSQL server compiled to WebAssembly inside your Go process. The binaries live here so that pgmem itself stays small, and so each Postgres major version can be its own Go module.

## Engines

| Module | PostgreSQL | Tag prefix |
|---|---|---|
| `github.com/feverenergy/pgmem-engine/pg17` | 17.11 | `pg17/` |

Each module embeds five files and exports them as `FS` (an `embed.FS`) along with a `Version` string. There is no other API. You normally don't import an engine directly.

## Choosing an engine from pgmem

pgmem requires `pg17` and runs Postgres 17 by default, so there is nothing to import. To name the engine yourself, set it in the options:

```go
srv, err := pgmem.Start(ctx, pgmem.Options{Engine: pgmem.Postgres17})

var tpl = pgtest.Template{Key: migrationsHash, Migrate: migrate, Engine: pgmem.Postgres17}
```

`pgmem.Engine` is a `Version` and an `FS`, so any module laid out like `pg17` works. A future `pg18` module would be used as `pgmem.Engine{Version: pg18.Version, FS: pg18.FS}`. One process can run both. pgmem keeps each engine's decompressed files, and pgtest its servers and cached images, apart by `Version`, so each engine needs a `Version` of its own. Pass the same `FS` value every time, for example from a package variable. `fs.Sub` returns a new value on each call, and pgmem rejects a second `FS` under a `Version` already in use.

## Tags

This repo has one Go module per Postgres major, so tags carry the directory as a prefix. A tag is `pg17/v0.<Postgres minor>.<rebuild>`. Postgres 17.11 is `pg17/v0.11.0`, a rebuild of 17.11 is `pg17/v0.11.1`, and 17.12 will be `pg17/v0.12.0`. Engine versions don't follow pgmem's version.

## How the files are built and verified

The files are produced by pgmem's build, with `make assets` in the pgmem repository. `pg17/PROVENANCE` lists the build inputs: the Postgres, wasi-sdk and Ubuntu snapshot pins, the Bison, flex and Go versions, the sha256 of each patch, and a sha256 over pgmem's build sources, with the rule for computing it.

`pg17/SHA256SUMS` has the sha256 of the three archives, and `pg17/INPUTS.SHA256SUMS` has the sha256 of each archive's decompressed contents. Both are embedded. `go test ./...` in `pg17` checks that each embedded archive matches `SHA256SUMS`, and pgmem's tests check the decompressed files against `INPUTS.SHA256SUMS`. pgmem's `make assets-check` rebuilds from source and fails, naming the file, if a file here has drifted from what the build produces.

## Licenses

The Go code is Apache-2.0, in `LICENSE`. pgmem's compat layer and Postgres patches, compiled into `postgres.wasm.zst`, are Apache-2.0 too. `pg17/` has a copy of `LICENSE` and a `NOTICE` with paths relative to it, so the module carries both. The compressed files contain PostgreSQL, wasi-libc and other third-party code. Their license texts are in `pg17/LICENSES/`, with an index in `pg17/LICENSES/README.md`.
