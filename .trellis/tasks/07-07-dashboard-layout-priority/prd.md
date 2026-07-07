# Dashboard Layout Priority

## Goal

Make the dashboard information flow match the user's work sequence: choose or search a product rule, inspect search results immediately, then adjust secondary configuration.

## Confirmed Facts

- The active Web panel lives in `template/html/index.html`.
- The four required areas are 商品规则, 搜索结果, 规则详情, and Telegram 通知.
- The previous desktop order placed notification or rule detail content in the center, forcing the user to look away from the rule/search flow.
- The current deployed commit is `39a05b5 Reorder dashboard columns`.

## Requirements

- On desktop, use the visual order 商品规则 -> 搜索结果 -> 规则详情 / Telegram 通知.
- Keep Telegram notification settings secondary on the right side, below rule details.
- Preserve mobile stacking without horizontal overflow.
- Do not change product search, Telegram, persistence, or Render configuration behavior.

## Acceptance Criteria

- [x] Desktop browser measurement shows grid areas `rules preview editor`.
- [x] Mobile browser measurement shows no horizontal overflow.
- [x] `node --check` passes for the inline script.
- [x] `go test ./...` passes.
- [x] Render deploy reaches `live`.
- [x] Production `/health` returns `200 {"status":"ok"}`.

## Out Of Scope

- Restoring deleted user data.
- Adding new notification providers.
- Rebuilding the whole visual style again.
