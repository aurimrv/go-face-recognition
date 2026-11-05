## Why

The CLI currently hardcodes a single input (`images/unknown.jpg`) and a single output (`images/result.jpg`), so analyzing a chosen photo or an entire album requires editing code or rerunning manually. Extending positional args to accept files and folders unlocks batch recognition without changing the recognition pipeline.

## What Changes

- Accept zero or more positional path arguments: no args keeps the current default; file args analyze those images; directory args walk recursively (including subfolders).
- Write annotated outputs beside each source as `result_<stem>.jpg`, except the no-args default which remains `images/result.jpg`.
- Resolve paths once into a job list, then classify/draw/save sequentially after a single person-load and recognizer init.
- Skip already-generated result images and non-JPEG files when walking directories; continue on per-image errors in batch mode.
- Document the three invocation forms in the README.

## Capabilities

### New Capabilities

- `batch-image-input`: Resolve CLI positional paths into classify/draw jobs (default unknown image, explicit files, recursive directories) and define output naming, walk filters, and batch error behavior.

### Modified Capabilities

- (none)

## Impact

- [`cmd/main.go`](cmd/main.go): replace hardcoded paths with job resolution and a processing loop.
- New small helper package/module for path → job resolution (testable pure logic).
- README usage/installation examples.
- No changes to person loading, embedding enrollment, or recognition models; drawer remains JPEG-only.
