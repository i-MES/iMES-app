# README

## 开发环境

- `git clone ...`
- `cd frontend`
- `yarn`
- `cd ..`
- `wails dev` or `wails build`

## Frontend Dev

### ESLint

- `wails dev` 时强制执行 ESLint：package.json 中的 `"dev": "eslint --ext .ts,.vue src && vite"`
- ESLint 的配置（`.eslintrc.js`）
  - 使用 `vue-eslint-parser` 做主解析器，检查 .vue 文件
  - 使用 `@typescript-eslint/parser` 做 `<script>` 检查，辅助主解析器
  - rules：单引号、行尾不用分号

## Backend Dev

### Test

- `go test ./...`
  - `go test ./backend/... -v`
    - `-v` 可以打印 `t.Log()` 和 `fmp.Println()` 内容，否则只打印 `t.Error()`
