# 面板交互与通知体验修正

## Goal

- 修正商品提醒面板的品牌、主题、规则管理、搜索反馈、商品状态与 Telegram 状态表达.
- 将推送记录直接展示在 Telegram 面板下方, 并让保存、关闭等本地交互立即响应.

## Confirmed Facts

- 前端集中在 `template/html/index.html`, 使用内联 CSS 与 JavaScript.
- 深色主题变量块被错误嵌套在未闭合的 `:root` 中, 导致手动深色和跟随系统失效.
- 当前推送记录通过 `openLogs()` 弹窗展示, 商品已推送状态统一使用 `.product-card.dim`.
- Render 使用 `master` 自动部署, 服务名为 `smzdm-for-go`, 健康检查路径为 `/health`.

## Requirements

- 左上角使用从什么值得买官方域名取得的官方图标资产, 同时用作页面 favicon.
- 右侧 Telegram 面板下方直接显示最近推送记录, 行样式与用户参考图一致, 不再依赖弹窗查看.
- 搜索中显示持续可见的加载动画和进度文案, 不呈现静止卡住状态.
- 商品价格与价格备注使用红色; 已推送或不可推送商品进一步降低视觉权重.
- 修复浅色、深色、跟随系统三种主题, 系统主题变化时自动同步.
- 顶栏 Telegram 图标根据通知真实启用状态显示灰色或 Telegram 蓝色.
- 每张商品规则卡片提供低干扰删除入口, 并保留确认步骤; 非商品系统规则不显示删除.
- 新建规则入口位于现有规则列表之后, 随列表滚动, 不固定在侧栏底部.
- 保存与面板关闭采用即时本地更新; 网络请求在后台完成, 失败时明确提示并回滚必要状态.

## Acceptance Criteria

- [x] Go tests和内联 JavaScript 语法检查通过.
- [x] 桌面与移动浏览器验证主题、删除、新建规则、搜索动画、商品状态和内联推送记录.
- [x] 保存与 Telegram 面板折叠操作无固定两秒等待, UI 在点击后立即更新.
- [x] 提交推送后 Render 部署状态为 `live`, `/health` 返回 `200 {"status":"ok"}`.

## Out Of Scope

- Unrelated product behavior changes.
