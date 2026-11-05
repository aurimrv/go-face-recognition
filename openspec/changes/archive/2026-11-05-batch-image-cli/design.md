## Context

`cmd/main.go` loads persons, initializes a `go-face` recognizer, enrolls known faces, then classifies a single hardcoded path (`images/unknown.jpg`) and writes `images/result.jpg`. Classify and draw already take paths; only input/output wiring is fixed. The `images/` tree already contains nested albums and some non-JPEG files. Drawer encoding/decoding is JPEG-only.

This is a teaching fork; prefer a small, readable change over a full CLI framework.

## Goals / Non-Goals

**Goals:**

- Positional CLI: no args (default), one or more files, and/or directories (recursive).
- Output naming: default → `images/result.jpg`; otherwise `result_<stem>.jpg` beside each source.
- One-time person load + recognizer init; sequential per-image classify/draw/save.
- Directory walks include only JPEG inputs and skip prior result images.
- Continue on per-image failure when processing multiple jobs; keep total timing.

**Non-Goals:**

- PNG/AVIF (or other format) support in the drawer.
- Concurrent recognition.
- Flag-based CLI (`cobra`, `-o`, `-threshold`, etc.).
- Changes to `persons/` loading or enrollment logic.
- Separate output root directory.

## Decisions

### 1. Positional args only (no flag framework)

**Choice:** Parse `os.Args[1:]` as filesystem paths; empty means default job.

**Why:** Matches the agreed UX (`go run ./cmd`, `... a.jpg b.jpg`, `... album/`) and keeps the teaching surface small.

**Alternatives considered:** `flag`/`cobra` with `--recursive` and `-o` — more extensible, heavier than needed now.

### 2. Pure job-resolver helper + thin main loop

**Choice:** Add a small package (e.g. `internal/input` or under `internal/usecases`) with pure functions that map args → `[]Job{In, Out}`. `main` loads recognizer once and loops.

**Why:** Path rules (default naming, walk filters, skip `result_*`) are easy to unit-test without dlib. Recognition use cases stay unchanged.

**Alternatives considered:** All logic inline in `main` (fastest, hardest to test); full batch use case wrapping classify+draw (more abstraction than needed).

### 3. Output naming

**Choice:**

- No args: `images/unknown.jpg` → `images/result.jpg` (preserve README/backward behavior).
- Explicit file or walked file: same directory, `result_<stem>.jpg` where stem is the source basename without extension.

**Why:** Explicit examples keep the historic default while making batch outputs unambiguous and colocated.

**Alternatives considered:** Always `result_<stem>.jpg` (breaks `result.jpg` default); dump all results into one out dir (worse for album trees).

### 4. Walk filters

**Choice:** When expanding a directory, include only `.jpg`/`.jpeg` (case-insensitive). Skip basename `result.jpg` and names matching `result_*.jpg`. Explicit file args are still accepted as given (caller responsibility) but expected to be JPEG for draw/save.

**Why:** Avoids reprocessing outputs on re-runs and skips formats the drawer cannot save. Nested albums under `images/` are first-class.

### 5. Error policy

**Choice:**

- Zero resolved jobs → fatal.
- Missing/invalid path argument → fatal (bad invocation).
- Default single-job mode (no args): classify/draw failure remains fatal.
- Two or more jobs (or any explicit-arg batch): log error for that image and continue; exit non-zero if any job failed.

**Why:** Albums should not abort on one bad photo; the historic single-image demo should still fail loudly.

### 6. Sequential processing

**Choice:** Process jobs one at a time on the shared recognizer.

**Why:** `go-face`/dlib thread-safety for one recognizer is unclear; sequential matches current latency expectations and simplifies error handling.

## Risks / Trade-offs

- **[Risk] Explicit non-JPEG file args still fail at draw/save** → Mitigation: document JPEG-only; walk already filters; surface clear per-job errors.
- **[Risk] Name collision if source is already `result_foo.jpg` and not skipped (explicit arg)** → Mitigation: walks skip result names; explicit args overwrite the computed output path intentionally.
- **[Risk] Large albums are slow (sequential + model cost)** → Mitigation: acceptable for teaching; concurrency deferred.
- **[Trade-off] No flags** → Threshold and paths stay constants in `main` until a follow-up.

## Migration Plan

- No data migration. Existing `images/result.jpg` behavior preserved for no-arg runs.
- Update README examples; existing Docker `go run cmd/main.go` keeps working.
- Rollback: revert the change; no schema or model changes.

## Open Questions

- None blocking implementation; threshold remains the current hardcoded `0.3` unless a later change exposes it.
