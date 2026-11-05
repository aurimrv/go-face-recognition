## 1. Job resolver

- [x] 1.1 Add a small package (e.g. `internal/input`) with a `Job` type (`In`, `Out`) and `Resolve(args []string) ([]Job, error)`
- [x] 1.2 Implement no-args default: `images/unknown.jpg` → `images/result.jpg`
- [x] 1.3 Implement file args: output `result_<stem>.jpg` in the source directory
- [x] 1.4 Implement directory args via recursive walk with `.jpg`/`.jpeg` filter and skip of `result.jpg` / `result_*.jpg`
- [x] 1.5 Fail when a path is missing/invalid or when zero jobs resolve
- [x] 1.6 Add unit tests covering default, files, nested walk, filters, and error cases

## 2. Main loop

- [x] 2.1 Refactor `cmd/main.go` to call `Resolve(os.Args[1:])` after keeping one-time person load + recognizer init
- [x] 2.2 Process jobs sequentially: classify → draw → save; print per-image face results
- [x] 2.3 Apply error policy: fatal for default single-job failures; continue + non-zero exit for multi-job / explicit-arg batches
- [x] 2.4 Preserve overall timing output; optionally note per-job failures summary

## 3. Docs and smoke check

- [x] 3.1 Update README Usage / Installation to document `go run ./cmd`, file args, and directory args
- [x] 3.2 Smoke-run default path and at least one multi-file or album path to confirm outputs land as specified
