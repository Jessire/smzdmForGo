# Implementation Plan

1. Extend `file.Config` and HTTP config mapping with `GlobalHotConfig` plus normalization helpers.
2. Add global feed scanning, time-window filtering, comment-floor filtering, author matching, merge/dedupe, and deterministic sorting in `smzdm/smzdm.go`.
3. Add focused backend tests for config persistence and pure candidate filtering.
4. Add compact dashboard controls and load/save normalization in `template/html/index.html`.
5. Run Go tests, JavaScript syntax validation, Docker build, and browser verification.
6. Commit the intentional source changes, push `master`, monitor Render until `live`, and verify `/health`.
