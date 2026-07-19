# Implementation Plan

1. Extend config/API persistence with independent hot and author keyword lists.
2. Replace old product-rule coupling with pure title-keyword matching.
3. Add reusable dashboard token editors for hot keywords, author nicknames, and author keywords.
4. Replace global-hot presets with positive integer inputs and preserve custom values end to end.
5. Add and run backend/UI tests plus browser interaction verification.
6. Obtain a valid PostgreSQL DSN, update Render, verify `/health/db`, deploy, then verify save survives restart.
