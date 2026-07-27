# Design

- 保留 Go 单文件面板架构, 仅在现有 HTML、CSS 与 JavaScript 中做定向升级.
- 把官方 favicon 作为本地静态资产提交, 避免运行时依赖外部站点和跨域加载.
- 规则卡片删除按钮仅在悬停、键盘聚焦或当前选中时提高可见度, 并阻止事件冒泡后调用现有确认对话框.
- 推送记录改为 Telegram 面板内的紧凑列表, 每行包含缩略图、截断标题、时间和状态; 顶栏图标直接读取 `enabled + botToken + chatId`.
- 搜索加载使用旋转环、脉冲轨道和省略号动画, 尊重 `prefers-reduced-motion`.
- 保存采用乐观状态: 立即清除脏状态并重绘, 请求失败时恢复脏状态并提示; Telegram 表单立即退出编辑态.
- 主题仅依赖 CSS 变量和 `data-theme`, 修正变量作用域并同步 `color-scheme`.

