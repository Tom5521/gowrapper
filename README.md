# gowrapper

Packs a directory into a self-extracting Go binary. That's it.

You point it at a folder — your app, its libraries, whatever it needs to run — and it spits out a single executable. When that executable runs, it decompresses everything into a temp directory, launches your binary, then cleans up. The end user never sees any of it.

```
your bundle/
├── myapp
├── libfoo.so
└── assets/
         ↓
   gowrapper
         ↓
  myapp-portable   ← ships everything, extracts on run, leaves no trace
```

---

## Why

Distributing apps with native dependencies is annoying. You either write an installer, ask users to install libraries themselves, or maintain a package for every distro. gowrapper is the lazy option: bundle the whole thing and ship one file. Works well for Qt/GTK apps, CLI tools with native deps, or anything compiled in a language that doesn't have Go's static linking story.

---

## Install

```bash
go install github.com/Tom5521/gowrapper@latest
```

Requires Go 1.21+ and a working Go toolchain (it calls `go build` internally).

---

## Usage

```
gowrapper --bundle <dir> --bin <path-inside-bundle> --out <output-file> [flags]
```

Three flags are required: `--bundle`, `--bin`, and `--out`. Everything else is optional.

```bash
gowrapper \
  --bundle ./dist/myapp \
  --bin    myapp \
  --out    myapp-portable
```

```bash
# Windows GUI app — no console window, lower compression for faster startup
gowrapper \
  --bundle      ./dist/myapp \
  --bin         myapp/myapp.exe \
  --out         myapp-portable.exe \
  --name        MyApp \
  --windowsgui \
  --compression-level 3
```

### All flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--bundle` | `-b` | required | Directory to pack (binary + all its dependencies) |
| `--bin` | | required | Path to the target binary *inside* the bundle |
| `--out` | `-o` | required | Output path for the generated executable |
| `--name` | `-n` | `""` | App name — used as the temp directory prefix |
| `--args` | `-a` | none | Default arguments passed to the bundled binary |
| `--compression-level` | `-c` | `9` | lz4 compression level, 0 (fast) to 9 (smallest) |
| `--verbose` | `-v` | false | Pass `-v` to the Go compiler |
| `--go-args` | `-g` | none | Extra arguments forwarded to `go build` |
| `--windowsgui` | | false | Adds `-H=windowsgui` to ldflags (hides the console window on Windows) |

The compression level trades startup time for bundle size. Level 9 gives the smallest binary; level 0 (`lz4.Fast`) extracts the quickest. For most apps, somewhere in the middle is fine.

---

## How it works

1. The bundle directory gets packed into a `tar` archive and compressed with lz4.
2. A small Go program (the "template") is written to a temp directory alongside the compressed bundle.
3. gowrapper runs `go build` on that template, injecting the app name, binary path, and default args via `-ldflags`.
4. The result is a standalone binary with the compressed bundle baked in via `//go:embed`.

When the generated binary is executed, it decompresses the bundle into `os.TempDir()`, runs the target binary (forwarding any arguments), and removes the temp directory on exit. On Windows, panics surface as native error dialogs via [zenity](https://github.com/ncruces/zenity) instead of crashing silently.

---

## Project layout

```
gowrapper/
├── main.go           # CLI (cobra), orchestrates the build process
├── template/
│   └── main.go       # The wrapper code embedded into every generated binary
└── util/
    ├── compress.go   # tar + lz4 compression
    └── decompress.go # tar extraction
```

---

## Dependencies

- [spf13/cobra](https://github.com/spf13/cobra) — CLI
- [pierrec/lz4](https://github.com/pierrec/lz4) — compression
- [ncruces/zenity](https://github.com/ncruces/zenity) — native error dialogs in the generated binary

---

## License

MIT
├── main.go           # CLI (cobra), orchestrates the build process
├── template/
│   └── main.go       # The wrapper code embedded into every generated binary
└── util/
    ├── compress.go   # tar + lz4 compression
    └── decompress.go # tar extraction

```

---

## Dependencies

- [spf13/cobra](https://github.com/spf13/cobra) — CLI
- [pierrec/lz4](https://github.com/pierrec/lz4) — compression
- [ncruces/zenity](https://github.com/ncruces/zenity) — native error dialogs in the generated binary

---

## License

[MIT](LICENSE)
