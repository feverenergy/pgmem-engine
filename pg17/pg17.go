// Package pg17 holds the compressed PostgreSQL 17 files that pgmem runs.
//
// The package has no code beyond the embedded files. FS has postgres.wasm.zst
// (the server compiled to WebAssembly), share.tar.zst (the share directory),
// cluster.tar.zst (an initialised data directory), SHA256SUMS, which lists
// the sha256 of each of the three archives, and INPUTS.SHA256SUMS, which
// lists the sha256 of each archive's decompressed contents. pgmem reads these
// files and decompresses them itself.
package pg17

import "embed"

// FS has the engine files at its root.
//
//go:embed postgres.wasm.zst share.tar.zst cluster.tar.zst SHA256SUMS INPUTS.SHA256SUMS
var FS embed.FS

// Version is the PostgreSQL release the files were built from.
const Version = "17.11"
