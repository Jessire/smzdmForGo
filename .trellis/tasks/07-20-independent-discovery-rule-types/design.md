# Technical Design

## Rule list model

The browser keeps `selectedRuleKind` (`product`, `hot`, or `author`) plus the selected product-rule index. `renderRules()` projects product rules and the two singleton discovery configurations into one list. Discovery cards are never deleted; their switch controls their existing enabled flag.

## Creation flow

The existing plus button opens a small three-choice dialog. Selecting 商品规则 appends a normal empty keyword rule. Selecting 全站热门 or 关注作者 selects that system card and marks its enabled flag true. This gives users an explicit rule type choice without duplicating persisted configuration.

## Editors and preview

The shared detail panel uses `data-editor-kind` to expose only the relevant fields:

- 商品规则: existing keyword, filter, price, and threshold fields.
- 全站热门: time window, comment floor, and independent hot title keywords.
- 关注作者: author nickname chips and independent author title keywords.

`/discoverySearch` accepts the existing global-hot request shape plus a type. It scans the recent full-site feed, applies the selected independent filter, and returns the same item payload as `/productSearch`, so the current preview renderer stays shared.

## Safety

Preview scans use a bounded page count. Production monitoring retains its existing broader scan. The read/write config contract remains `globalHot`, so saved PostgreSQL rows stay compatible.
