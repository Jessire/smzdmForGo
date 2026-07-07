# Design

## Architecture

The rebuild is CSS-first. Existing DOM IDs and JavaScript event bindings remain intact:

- `#rulesBody`, `#searchAllRules`, `#addRule`, `#reloadProductConfig`
- `#previewHero`, `#previewCards`, `#openSearchPage`
- rule editor inputs and `#saveProductConfig`
- Telegram inputs and `#saveNotifyConfig`

The implementation adds a final `trellis-ui` design layer after the previous CSS. This preserves behavior while replacing the visible hierarchy.

## Layout

Desktop:

```text
topbar: 商品提醒规则 / 保存状态 / TG 状态 / 主题状态

商品规则 | 规则详情      | Telegram 通知
商品规则 | 搜索结果      | 搜索结果
```

All desktop regions sit inside one shared app canvas. Functional areas remain visible, but separation is expressed with soft gutters, low-contrast borders, and aligned panel edges rather than heavy split lines or black gaps.

Mobile:

```text
商品规则
规则详情
搜索结果
Telegram 通知
```

## Visual System

- Dense admin surface with restrained color.
- Shared app canvas with quiet internal panels, low shadow, and consistent 8px-radius geometry.
- Icon colors use semantic categories only: blue for search/save/open primary actions, green for enabled/success, red for delete/destructive, neutral for field icons and secondary metadata.
- Product and rule images stay in color.
- Rules occupy the full left column and scroll internally.
- Rule editor and Telegram share the right workspace top row at about 70/30.
- Search preview spans the full right workspace below the editor and Telegram.
- Search results render as compact independent product cards on desktop. The preview area sizes to actual result count, shows up to six cards before internal scrolling, and does not reserve empty result slots.
- Empty search state spans the full preview area instead of occupying one card slot.

## Risks

- The page has multiple historical CSS blocks, so the final CSS block must use enough specificity to override previous `.zero-redesign` rules.
- Body class must keep `zero-redesign` until legacy CSS/JS dependencies are removed in a separate cleanup.
