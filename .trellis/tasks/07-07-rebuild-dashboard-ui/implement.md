# Implementation Plan

1. Mark task in progress.
2. Read frontend/shared Trellis specs.
3. Preserve `trellis-ui` body class.
4. Update the final `Trellis UI rebuild` CSS/JS layer to match the integrated canvas, inline metrics, compact Telegram, and footer action placement decisions.
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
  - Desktop grid areas: topbar, rules, editor, notify, preview, all reading as one shared app canvas.
  - Search result preview: independent product cards, sized to actual result count, up to 6 before internal scrolling.
  - Filter chips: fixed-height box with internal scroll when chips exceed two rows.
  - Rule detail: no enable switch; footer action group contains preview current rule and save-and-enable.
  - Telegram: title-bar plane enable control; body contains only Bot Token, Chat ID, Parse Mode / HTML, and link preview.
  - Icon colors: blue primary/search/open, green enabled/success, red destructive, neutral field/metadata.
  - Mobile: `scrollWidth <= clientWidth`.
