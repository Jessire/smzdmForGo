# 独立发现规则类型

## Goal

- Make 全站热门 and 关注作者 first-class entries in 商品规则 instead of hidden advanced settings.
- Let the add button choose 商品规则, 全站热门, or 关注作者.
- Let each rule type render its own editor and use the shared 搜索预览 panel.

## Confirmed Facts

- `GlobalHotConfig` already stores the independent hot and author filters in `file/fileUtis.go`.
- The dashboard currently renders only `keywordRules` in `renderRules()` and keeps discovery settings in `#advancedRuleFields` in `template/html/index.html`.
- `/productSearch` only supports keyword search in `route_search.go`; discovery preview requires a dedicated full-site search path.

## Requirements

- Move the full-station hot and followed-author controls out of the advanced settings block.
- Render them as singleton system rule cards alongside ordinary product rules, with independent enable toggles and type-specific metadata.
- The `+` action presents 商品规则 / 全站热门 / 关注作者. Product adds a normal editable rule. The other choices select and enable their system rule without creating duplicates.
- Selecting a system rule shows only its relevant fields in the rule detail panel.
- Searching a system rule populates the existing 搜索预览 with full-site matches under that rule's filters.
- Keep normal product-rule search, saved configuration, Telegram, and push behavior unchanged.

## Acceptance Criteria

- [ ] Go tests cover discovery preview matching and configuration persistence.
- [ ] Dashboard script parsing and UI marker tests pass.
- [ ] Browser verification covers add-type selection, rule card selection, each editor, and search preview at desktop and mobile widths.

## Out Of Scope

- Multiple independent instances of 全站热门 or 关注作者.
- Changes to Telegram delivery, existing product rule matching, or the layout outside the specified rule-management flow.
