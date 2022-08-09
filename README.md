# README

## 基础说明

- 主色调：中国红(`#A7535A`)、微软蓝(`#1c7bc9`)

## 开发环境

- install [node.js](http://nodejs.cn/)、[go](https://go.dev/doc/install)
- install [wails](https://wails.io/docs/gettingstarted/installation)
  - `go get github.com/wailsapp/wails/v2@v2.0.0-beta.38`
- `npm config set registry http://registry.npmmirror.com`
- `git clone ...`
- `cd frontend`
- `yarn`
- `cd ..`
- `wails dev` or `wails build`

### 升级 wails 版本

- 如果已安装老版本的 wails，建议先卸载：
  - `` rm `go env GOPATH`/bin/wails ``
  - `` rm `go env GOPATH`/pkg/mod/github.com/wailsapp -rf ``
- 如果新版本安装不成功，可以进入目录手动安装：
  - `` cd `go env GOPATH`/pkg/mod/github.com/wailsapp/wails/v2@<已下载的新版本> ``
  - `cd cmd/wails`
  - `go install`: 生成 wails 可执行程序

### 执行测试用例

```sh
$ go test ./backend/... [-count n] [-v]
```

或

```sh
$ go test -bench <XXX> ./backend/...
```

`-v` 可以查看 `t.Log()` 输出的内容。

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

### Run Model

1. `wails dev`: 对源代码文件实时监控，修改后自动重新 `wails dev`，可用于日常开发。
2. `F5`,`Shift-F5`: 不对源代码文件实时监控，首先 `go build --debug...`，然后执行 `build/bin/iMES-app`，可以设置断点，对源代码文件的修改会导致断点定位失效，可用于单步 debug 定位故障。

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

1. `wails build`
2. `cd frontend && yarn build`

两种方法都会重新生成 `frontend/dist` 路径及文件，`wails dev` 不会重新生成。

### How to view vuetify docs locally?

- `git clone https://github.com/vuetifyjs/vuetify`
- `cd vuetify`
- `yarn`
- `yarn build`（支持分别编译：`yarn build dev`,`yarn build api`,`yarn build docs`）
- `http-server packages/docs/dist -p 8080`

### How to view wails doc locally?

- `git clone https://github.com/wailsapp/wails.git`
- `cd wails/website`
- `yarn build && http-server build -p 8081` or `yarn start`

### How to view godoc locally?

- `cd iMES-app`
- `godoc -http=localhost:8082 -index`
  - godoc 工作需要 go.mod，是针对某个 go module 生成 doc，且生成其依赖包指定版本 doc，而不是生成某个依赖包的最新版本的 doc —— 所以 godoc 需要在 app 内执行。
  - `-index` 可以生成索引以便搜索。
- if error: `//go:build comment without // +build comment`
  - check godoc version, it is one cmd of `golang.org/x/tools`, not a standard library.
  - 就像上面重新编译 wails 一样，可以进入新的 godoc 路径，重新编译，如：
  - `` cd `go env GOPATH`/pkg/mod/golang.org/x/tools@<新版本>/cmd/godoc ``
  - `go build -o ../../../../../../../bin/`
  - 再执行 `godoc` 时可查看网页的 footer 是否显示新的版本号。
  - 原因：
    - ` +build` 是老版本的 go 编译指令，`//go:build` 是 go1.17 之后新的。
    - 为了兼容性，老版本的 godoc 要求两个都要写，发现只有 1 个就报以上错误。
    - 但新版本的 godoc 已经能够更智能的处理，不再做这个要求。

### `wails dev` 因端口不对无 UI

```
DEB | [DevWebServer] Waiting for frontend DevServer 'http://localhost:3000' to be ready
   vite v2.9.13 dev server running at:

  > Local: http://localhost:3001/                                                                                                                                                                                                                                    14:54:36
  > Network: use --host to expose
```

因为 wails 前端使用 3000 端口，上面提示说明 3000 被占用，又可以是因为上一个 `wails dev` 没有正常关闭。

```sh
➜  iMES-app git:(alpha) ✗ netstat -ap|grep 3000
(Not all processes could be identified, non-owned process info
 will not be shown, you would have to be root to see it all.)
tcp        8      0 localhost:3000          0.0.0.0:*               LISTEN      85879/node
tcp        1      0 localhost:3000          localhost:33848         CLOSE_WAIT  85879/node
tcp        1      0 localhost:3000          localhost:33912         CLOSE_WAIT  -
tcp      422      0 localhost:3000          localhost:33910         CLOSE_WAIT  -
```

说明 3000 被 85879 的 node 占用。

```sh
ps -ef|grep 85879
me         85879    1774  0 11:34 ?        00:00:23 /data/software/node/node-v18.2.0-linux-x64/bin/node /data/kevin/workspace/kproject/imes/iMES-app/frontend/node_modules/.bin/vite
me         85886   85879  0 11:34 ?        00:00:12 /data/kevin/workspace/kproject/imes/iMES-app/frontend/node_modules/esbuild-linux-64/bin/esbuild --service=0.14.43 --ping
```

进一步确认是 iMES-app 的 vite 导致 node 没有正常关闭。所以：

```sh
kill -9 85879
```

或者

```sh
$ pkill node
```

## Others

### git 提交规范

遵守 [约定式提交规范](https://www.conventionalcommits.org/zh-hans/)

- feat：新功能（feature）
- fix：修补 bug
- refactor：重构（即不是新增功能，也不是修改 bug 的代码变动）
- docs：文档（documentation）
- style： 格式（不影响代码运行的变动）
- test：增加测试
- chore：构建过程或辅助工具的变动
