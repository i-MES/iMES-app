# Slicksort

https://vue-slicksort.netlify.app/

## 源码解析

- 上层 3 个组件 -- 最终都是生成 `<div>`
  - SlickList —— 基于 ContainerMixin
  - SlickItem —— 基于 ElementMixin
  - DragHandle
- 底层 2 个组件
  - ContainerMixin
  - ElementMixin
- 最底层
  - Manager —— 2 个 Mixin 都使用
- 其他
  - 作为 plugin 时需要 `app.use(Slicksort)，此时会初始化
    - HandleDirective：添加 handle 指令
    - SlicksortHub
