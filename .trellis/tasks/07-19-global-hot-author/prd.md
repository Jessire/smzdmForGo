# 全站热门与关注作者推送

## Goal

- Add an optional full-site hot-deal channel that does not require keywords.
- Let the user choose a 3, 6, or 12 hour publication window and a 100+ or 200+ comment floor.
- Let the user optionally apply the existing product keyword rules as a second filter.
- Let the user optionally follow exact SMZDM author nicknames and push their new good-price posts.

## Confirmed Facts

- `smzdm.GetGoods(page, "")` already retrieves the full good-price feed in `smzdm/smzdm.go`.
- The feed exposes Unix publication time as `publish_date_lt`, cumulative comments as `article_comment`, and author nickname as `article_referrals`.
- The feed has no usable comment-sort or author-filter parameter; filtering and comment sorting must happen locally.
- Existing keyword rules, price filters, value/comment thresholds, Telegram push, and article-id dedupe must remain compatible.
- Product configuration is persisted as a JSON blob through `db/settings.go`, so these fields do not require a relational schema migration.

## Requirements

- Add a persisted `globalHot` configuration with `enabled`, `windowHours`, `minCommentNum`, `applyKeywordRules`, `followAuthorsEnabled`, and `followedAuthors`.
- When global hot is enabled, scan all feed pages needed to cover the selected publication window, retain items meeting the comment floor, optionally apply existing keyword rules, sort by comments descending and publication time descending, then merge with existing results and deduplicate by article id.
- When author following is enabled, exact-match normalized `article_referrals` against configured nicknames and include matching good-price posts from the same global scan without requiring the hot comment floor.
- Preserve existing keyword-only behavior when both new channels are disabled.
- Add compact panel controls for both channels without changing the existing 商品规则, 搜索预览, 规则详情, or Telegram notification layout hierarchy.
- Validate configuration bounds server-side: window is 3, 6, or 12 hours; comment floor is at least 100; author names are trimmed and deduplicated.

## Acceptance Criteria

- [ ] Existing Go tests pass and new unit tests cover global hot filtering, author matching, sorting, and config round-trip.
- [ ] Inline dashboard JavaScript passes `node --check`.
- [ ] Real browser verification confirms controls load, save, reload, and do not disturb the existing four-panel layout.
- [ ] Docker image builds successfully.
- [ ] Render reports the deployed commit as `live` and `/health` returns HTTP 200 with `{"status":"ok"}`.

## Out Of Scope

- Tracking comments added during the window rather than cumulative comments on recently published items.
- Following all SMZDM content types outside the current `good_price` feed.
- Replacing the existing push-history storage in this task.
