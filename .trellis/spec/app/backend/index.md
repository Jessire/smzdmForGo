# Backend Spec

## Scope

Applies to `main.go`, `route*.go`, `db/`, `file/`, `smzdm/`, and `push/`.

## Pre-Development Checklist

- Read `.trellis/spec/guides/index.md`.
- Preserve HTTP handler registration in `main.go` unless changing public routes intentionally.
- Run Go tests after backend changes.

## Local Patterns

- `main.go` owns service startup, route registration, current config synchronization, and the product cron loop.
- Product config API shape is defined in `route_config.go` through `productConfigRequest`, `telegramConfigRequest`, and `keywordRuleConfigRequest`.
- Request cleanup should use existing helpers in `route_config.go`, including `cleanWords`, `nonNegativeInt`, `nonNegativeFloat`, and `normalizedParseMode`.
- Search preview is handled by `ProductSearchHandler` in `route_search.go`; it accepts a single rule and returns `keyword`, `openUrl`, and product `items`.
- Runtime config is stored through `db.NewDB(userDbPath)` and loaded over the default YAML config during startup.

## Avoid

- Do not read Telegram credentials from new environment variables for normal user configuration; the panel and persisted config are the current source of truth.
- Do not mutate global config without `setCurrentConfig`.
- Do not add route behavior that bypasses existing JSON response helpers.
