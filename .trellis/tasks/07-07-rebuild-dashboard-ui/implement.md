# Implementation Plan

1. Mark task in progress.
2. Read frontend/shared Trellis specs.
3. Add `trellis-ui` body class.
4. Add a final `Trellis UI rebuild` CSS block that replaces desktop and mobile layout presentation.
5. Verify script syntax with `node --check`.
6. Run Go tests.
7. Start local server and inspect layout with Agent Browser and Chrome at desktop and mobile widths.
8. Commit only `template/html/index.html` and Trellis task files; leave `data/users.db` uncommitted.

## Validation Commands

- `node --check` on the inline script extracted from `template/html/index.html`.
- `C:\Users\jery3\.codex\tools\go1.26.4\go\bin\go.exe test ./...`
- `git diff --check`
- Agent Browser DOM checks for desktop and mobile widths.
- Chrome render checks at 1440px and 390px:
  - Desktop grid areas: topbar, rules, editor, notify, preview.
  - Search result preview: 6 product cards, 3 columns x 2 rows.
  - Filter chips: fixed-height box with internal scroll when chips exceed two rows.
  - Mobile: `scrollWidth <= clientWidth`.
