# Technical Design

## Persistence

The existing PostgreSQL JSON settings blob remains authoritative. The current loss is caused by an unreachable Aiven DSN, not JSON serialization. Keep `REQUIRE_DATABASE_URL=true`; once a valid DSN is supplied, verify save, restart, and reload against the same database row.

After Aiven is running, GitHub Actions calls the database health endpoint every 10 minutes so Render wakes before the request and the endpoint performs a real PostgreSQL ping. The running application performs the same lightweight ping every 10 minutes while its process is active.

## Independent Filters

Extend `GlobalHotConfig` with `HotKeywords` and `AuthorKeywords`. `ApplyKeywordRules` remains decodable for backward compatibility but is no longer used by matching or shown in the UI.

- Global hot: time window -> comment floor -> optional `HotKeywords` title match.
- Followed author: exact author nickname -> optional `AuthorKeywords` title match.

Each list uses OR semantics and case-insensitive substring matching against the product title.

The time window and comment threshold accept any positive integer. Non-positive or missing values normalize to the existing defaults of 3 hours and 200 comments.

## UI

Use a reusable token-editor helper for three fields: followed author nicknames, global hot keywords, and author keywords. Reuse the existing token appearance without the common-filter pin control. Enter, comma, Chinese comma, Tab, multi-value paste, Backspace, and remove actions behave consistently.
