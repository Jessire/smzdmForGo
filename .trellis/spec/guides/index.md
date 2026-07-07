# Shared Guides

## Pre-Development Checklist

- Read the active task `prd.md` before editing.
- Check `git status --short` and avoid committing runtime data.
- For UI changes, inspect `template/html/index.html` and verify layout with Agent Browser or equivalent real browser tooling.
- For deployment changes, read `render.yaml`, `Dockerfile`, and `readme.md`.

## Project-Wide Rules

- Keep user-facing output and UI labels in Chinese unless preserving protocol or API field names.
- Do not commit `data/users.db`; it is runtime data and local browser/server runs can mutate it.
- Go verification command: `C:\Users\jery3\.codex\tools\go1.26.4\go\bin\go.exe test ./...`.
- HTML inline script verification: extract the `<script>` block from `template/html/index.html` and run `node --check`.
- Deployment verification requires Render status `live` and `https://smzdm-for-go.onrender.com/health` returning `200 {"status":"ok"}`.
