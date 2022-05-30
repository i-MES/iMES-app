# README

## 开发环境

- `wails dev`
- `cd frontend` - `npm run dev`
- `wails build` to `build/bin/`

## Dev

- `go test ./...`
  - `go test ./backend/... -v`
    - `-v` 可以打印 `t.Log()` 和 `fmp.Println()` 内容，否则只打印 `t.Error()`
