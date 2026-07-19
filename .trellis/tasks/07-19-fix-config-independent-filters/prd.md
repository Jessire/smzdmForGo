# 修复配置持久化与独立筛选

## Goal

- Preserve saved production configuration across deployments.
- Give 全站热门 and 关注作者 their own independent keyword filters.
- Render author names and both keyword lists as removable token chips matching the existing filter-word interaction.

## Confirmed Facts

- Production `/health/db` currently returns `db_down` because the configured Aiven hostname does not resolve.
- Production `/productConfig` returns the repository defaults, confirming database-backed configuration is not being loaded.
- The current global-hot implementation reuses existing product rules through `ApplyKeywordRules`; the user explicitly rejects that coupling.
- The current author input is a plain text field and serializes only `followedAuthors`.

## Requirements

- Add `hotKeywords` and `authorKeywords` to `GlobalHotConfig`, API mapping, JSON persistence, and dashboard normalization.
- Global hot items must match any `hotKeywords` entry when the list is non-empty; otherwise the comment threshold alone is sufficient.
- Followed-author items must first match an exact author nickname, then match any `authorKeywords` entry when that list is non-empty.
- Remove the UI and runtime dependency on existing product rules for both new channels.
- Replace the plain author input with token chips, and add equivalent token editors for hot and author keywords.
- Replace the global-hot preset selects with positive integer inputs so the time window and minimum comment count can be entered freely.
- Do not report a successful configuration save when PostgreSQL is unavailable.
- Do not replace the broken PostgreSQL dependency with ephemeral Render-local SQLite.
- Keep the Render service and Aiven database warm by probing `/health/db` every 10 minutes from GitHub Actions and from the running application.

## Acceptance Criteria

- [x] Go tests cover independent hot keywords, exact author matching, author keywords, and config round-trip.
- [x] Dashboard JavaScript syntax check passes.
- [x] Browser verification confirms Enter/paste creates chips and removal updates saved JSON.
- [x] Browser verification confirms custom time/comment values survive save and reload without horizontal overflow at desktop and mobile widths.
- [x] Production deployment occurs only after a valid PostgreSQL DSN is available and `/health/db` returns 200.

## Out Of Scope

- Provisioning or choosing a new database without user authorization.
- Following non-`good_price` content types.
