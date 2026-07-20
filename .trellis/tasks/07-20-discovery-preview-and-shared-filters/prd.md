# 发现规则预览与通用筛选修正

## Goal

- 将新增规则选择器贴近左侧加号, 保持规则管理上下文.
- 让全站热门和关注作者复用商品规则的通用筛选条件, 另加发现时间范围.
- 修复发现预览结果展示, 包括作者信息和空结果的可解释性.

## Requirements

- 新增规则菜单锚定在左侧加号附近, 不居中覆盖全屏.
- 全站热门和关注作者详情不显示重复启用开关.
- 通用筛选字段在三类规则中保持一致; 发现规则额外显示时间范围和各自关键词/作者字段.
- 发现预览使用有限全站扫描, 正确返回符合条件的商品和作者字段.
- 不改变 Telegram 配置和普通商品规则的推送语义.

## Acceptance Criteria

- [ ] Go tests and inline JavaScript checks pass.
- [ ] Desktop and mobile browser verification covers the anchored picker, shared fields, and discovery preview.
- [ ] Production health checks pass after deployment.
