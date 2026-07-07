# Deployment Spec

## Scope

Applies to `Dockerfile`, `render.yaml`, README deployment notes, and Render verification.

## Pre-Development Checklist

- Read `.trellis/spec/guides/index.md`.
- Read `C:\Users\jery3\.codex\skills\docker-deploy-runbook\SKILL.md` before Render deployment work.
- Verify `/health` after deployment.

## Local Patterns

- Render service is a Docker Web Service with health check path `/health`.
- Production should use PostgreSQL through `DATABASE_URL`.
- `REQUIRE_DATABASE_URL=true` prevents production fallback to local SQLite.
- Render deployment verification must include deployed commit, deploy status, and health response.

## Avoid

- Do not commit production secrets.
- Do not rely on container-local SQLite for production state.
- Do not treat a pushed commit as deployed until Render reports `live`.
