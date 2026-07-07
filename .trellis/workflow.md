# Trellis Workflow

## Phase Index

1. Planning
   - Capture the user-facing goal in `prd.md`.
   - Inspect repository evidence before asking product questions.
   - Use `design.md` and `implement.md` only when the change has meaningful technical risk or multiple implementation steps.

2. Development
   - Read the active task artifacts.
   - Read `.trellis/spec/guides/index.md` and each relevant package or layer index.
   - Modify only files required by the task.

3. Verification
   - Run the commands listed in the task acceptance criteria.
   - For UI changes, verify the real browser layout at desktop and mobile widths.
   - For deployment changes, verify Render health checks after rollout.

4. Completion
   - Commit only intentional source and Trellis files.
   - Do not commit runtime data such as `data/users.db`.
   - Record unresolved risks in the final response or task notes.

## Request Triage

- Use a lightweight Trellis task for focused UI, config, or route changes.
- Use full planning artifacts for new product features, persistence changes, deployment architecture, or notification behavior changes.
- Use direct edits only for trivial typo or copy changes when the user does not need task traceability.

## Current Project Notes

This repository is a Go Web Service for smzdm product monitoring with a single static HTML panel in `template/html/index.html`, Go HTTP routes in `route*.go`, and persistence under `db/`.
