# Licenses

The Go code in this module is Apache-2.0, in `../LICENSE`. The compressed files inside it contain third-party code, and this directory has the license texts for it. Texts are copied verbatim from the sources in the last column.

The files cover `postgres.wasm.zst`, `share.tar.zst` and `cluster.tar.zst`.

| Component | Version | License | Where used | License file (source of the text) |
|---|---|---|---|---|
| pgmem's compat layer and PostgreSQL patches | as in pgmem's build | Apache-2.0 (Fever Energy AS) | Compiled into `postgres.wasm` | `../LICENSE`, `../NOTICE` |
| PostgreSQL | 17.11 | PostgreSQL License | The server, `btree_gist`, the regression C libraries, the conversion procs, the initdb'd cluster image | `postgresql.txt` (`COPYRIGHT` in `postgresql-17.11.tar.bz2`) |
| PostgreSQL bundled code: sha2, md5/sha1, snprintf, strftime, qsort/bsearch, mkdir -p, wide-char compare, inet, timingsafe_bcmp, the regex package (Spencer, Tcl), GEQO (Genitor), HyperLogLog, fasthash, Ryu | as in PostgreSQL 17.11 | BSD-3-Clause, ISC, Spencer and Tcl, Genitor, MIT, Boost-1.0 | Compiled into `postgres.wasm` | `postgresql-third-party.txt` (each file's notice, verbatim from `postgresql-17.11.tar.bz2`) |
| Snowball stemmers and stopword lists | v2.2.0, as vendored in PostgreSQL 17.11 | BSD-3-Clause | The stemmers in `postgres.wasm`, and the stopword files in `share.tar.zst` | `snowball.txt` (`COPYING`, snowballstem/snowball tag v2.2.0) |
| wasi-libc | wasi-sdk 34 (wasi-libc `2e6fb9d8ee0c`) | Apache-2.0 WITH LLVM-exception, Apache-2.0 and MIT (multi-licensed) | libc and the WASI bottom half linked into `postgres.wasm` | `wasi-libc-NOTICE.txt` (`LICENSE`), `wasi-libc-Apache-2.0-LLVM-exception.txt` (`LICENSE-APACHE-LLVM`), `wasi-libc-Apache-2.0.txt` (`LICENSE-APACHE`), `wasi-libc-MIT.txt` (`LICENSE-MIT`) |
| musl libc (inside wasi-libc) | as vendored by wasi-libc | MIT, with other permissive notices listed in the file | string, math, stdio and stdlib functions | `musl.txt` (`libc-top-half/musl/COPYRIGHT`) |
| cloudlibc (inside wasi-libc) | as vendored by wasi-libc | BSD-2-Clause | The WASI syscall wrappers | `cloudlibc.txt` (`libc-bottom-half/cloudlibc/LICENSE`) |
| dlmalloc (inside wasi-libc) | 2.8.6 | CC0 (public domain dedication) | `malloc`, `free` | `dlmalloc.txt` (the notice at the top of `dlmalloc/src/malloc.c`) |
| compiler-rt builtins (wasi-sdk 34, LLVM 23.1.0) | 23.1.0 | Apache-2.0 WITH LLVM-exception | `__multi3`, `__udivti3` and the other integer helpers | `compiler-rt.txt` (`compiler-rt/LICENSE.TXT`) |
| tzdata | the one in the PostgreSQL 17.11 tarball | Public domain | `share.tar.zst`: the timezone files compiled with `zic` | No license text, because the IANA time zone database is in the public domain |

## How the PostgreSQL notices were chosen

To choose the bundled PostgreSQL code, we ran `grep -rlE "Redistributions? (in|of) binary form|appear in all copies|retained in all copies|copyright notice|public domain"` over the 17.11 tree's `src` and `contrib`. We kept a candidate only if its functions are in the name section of `postgres.wasm` as pgmem's build produces it. These are present:

- `pg_sha224_init`, `pg_sha256_init`, `pg_md5_init`, `pg_sha1_init`;
- `pg_snprintf`, `pg_strftime`, `pg_qsort`, `bsearch_arg`, `pg_mkdir_p`, `pg_char_and_wchar_strncmp`;
- `pg_inet_net_ntop`, `pg_inet_net_pton`, `pg_inet_cidr_ntop`;
- `timingsafe_bcmp`, which is PostgreSQL's, because wasi-libc does not define it;
- `pg_regcomp` and `cclasscvec`, because `regcomp.c` `#include`s `regc_locale.c`;
- `gimme_edge_table`, `geqo_selection`, `geqo_copy`, `initHyperLogLog`;
- `spcache_insert` in `namespace.c`, which inlines fasthash;
- `double_to_shortest_decimal_bufn`.

`wstrcmp.c` is compiled and linked, but the linker drops its only function. Its notice is included anyway.

We left these out because they are not linked:

- `getopt`, `getopt_long`, `strlcpy` and `strlcat`, which are wasi-libc's and covered by `musl.txt`;
- `inet_aton` and `mkdtemp`;
- the other GEQO crossover files, because `cx`, `px`, `pmx`, `ox1`, `ox2` and `geqo_mutation` are compiled out;
- all of contrib except `btree_gist` and `spi`, neither of which has its own notice;
- `src/bin`, `src/interfaces` and `src/tools`;
- `src/backend/jit/llvm/SectionMemoryManager.LICENSE`, because the build uses `--without-llvm`.

The time zone code (`localtime.c`, `private.h`, `tzfile.h`, `zic.c`) and Markus Kuhn's `wcwidth` in `src/common/wchar.c` are in the public domain, so they need no notice. initdb is not shipped, because `initdb.wasm` is not embedded, so the files only it links need no notice here.

`postgres.wasm` links these parts of wasi-sdk:

- `libc.a`, which is wasi-libc with the musl top half, the cloudlibc bottom half and dlmalloc;
- `libsetjmp`, `libwasi-emulated-process-clocks` and `libwasi-emulated-mman`, as pgmem's `build/wasm-env.sh` sets;
- `libclang_rt.builtins.a`.

The binary has dlmalloc and compiler-rt symbols (`__multi3`, `__udivti3`, `__ashlti3`) and no `emmalloc` or `fts_open`. So emmalloc and musl-fts, which wasi-libc's LICENSE also lists, are not linked. Postgres is C, so `postgres.wasm` does not link libc++ or libc++abi, and their license is not included.

wasi-sdk's compiler and scripts, and the Ubuntu packages in the build image, are build tools. They are not part of the files beyond the libraries above.
