# README

## 开发环境

- install [node.js](http://nodejs.cn/)、[go](https://go.dev/doc/install)、[wails](https://wails.io/docs/gettingstarted/installation)
- `npm config set registry http://registry.npmmirror.com`
- `git clone ...`
- `cd frontend`
- `yarn`
- `cd ..`
- `wails dev` or `wails build`

### 切换 python 版本

```sh
$ cd backend/python
$ ./change_py_pc.sh 3.y.z
```

- 脚本会自动搜索本机通过 pyenv 安装的，且有共享库的 python 版本，y.z 填写不对时会列出。
  - 通过 `env PYTHON_CONFIGURE_OPTS="--enable-shared" pyenv install -v 3.y.z` 从源码编译安装 python 共享库
- 最后需拷贝 `LD_LIBRARY_PATH=...` 语句到 shell 窗口，执行后在运行 `wails dev` 或 `wails build`

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

## FAQ

### `frontend/dist` 目录下空

```sh
$ wails dev
......
ERROR:
main.go:14:12: pattern frontend/dist: cannot embed directory frontend/dist: contains no embeddable files
```

解决方法：

```sh
$ cd frontend/dist
$ yarn build
```
