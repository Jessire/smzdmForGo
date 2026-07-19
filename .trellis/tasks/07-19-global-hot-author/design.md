# Technical Design

## Data Flow

`productConfig` JSON -> `file.Config.GlobalHot` -> `smzdm.GetSatisfiedGoods` -> global feed scan -> local filters and sort -> existing Telegram pusher.

The scan uses the existing `/v1/list` request with an empty keyword and `order=time`. It keeps pages until the selected cutoff is covered, with a bounded page limit and a consecutive-empty-page stop condition to tolerate occasional out-of-order feed rows.

## Configuration

Add a `GlobalHotConfig` value to `file.Config`. The HTTP request/response structs mirror it. The existing JSON settings blob persists it without a database migration. Old saved configurations decode to disabled zero values and keep the current keyword path unchanged.

## Matching

- Hot candidates: publication timestamp inside the window, cumulative comments at or above the configured floor, not already pushed, and optionally matching at least one enabled keyword rule.
- Followed-author candidates: publication timestamp inside the scan window, exact normalized `article_referrals` match, and not already pushed. They bypass the hot comment floor but remain within the current `good_price` feed.
- Merge all sources by `ArticleId`, then sort comments descending and publication time descending before Telegram batching.

## UI

Place a compact discovery block inside the existing advanced rule area. It contains independent toggles for 全站热门 and 关注作者, selects for 3/6/12 hours and 100+/200+, a second-filter toggle, and an author nickname token input. Existing layout sections and Telegram controls remain untouched.

## Risks

The endpoint is an unauthenticated site API rather than the documented signed open API, so rate limiting remains an operational risk. A 12-hour full-site scan can require dozens of 100-row pages; requests are throttled and bounded.
