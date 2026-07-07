# Rebuild Dashboard UI

## Goal

Rebuild the Web panel UI around the product's real job: create scheduled product reminder rules, briefly verify a rule with search preview, then let Telegram scheduled push deliver reminders.

## Background

- The active UI is implemented in `template/html/index.html`.
- Existing backend routes and JSON contracts already support product config, search preview, image proxy, and Telegram config.
- The existing DOM has four functional areas: 商品规则, 搜索结果, 规则详情, and Telegram 通知.
- The user rejected patch-style visual tweaks and wants a real redesign.
- Product intent is a reminder-rule configuration tool. Search is temporary validation, not the primary product surface.
- Telegram is the delivery channel for scheduled reminders. It must be visible, but it does not deserve a full-height dedicated column.
- UI/Taste skill read: this is dense dashboard/admin product UI, not a marketing page. Use Taste only for anti-template/aesthetic constraints; use UI/UX Pro Max for forms, layout, accessibility, and interaction quality.
- UI/UX Pro Max direction: data-dense dashboard, low motion, high practical density, visible labels, submit feedback, accessible icon buttons, and light/dark support.

## Requirements

- Preserve existing DOM IDs, JavaScript behavior, route contracts, and feature coverage unless a later implementation note proves a wrapper change is required.
- Desktop uses a thin global top bar, about 56-64px, for product name, save status, Telegram health, and theme state only. It must not become a hero/header block.
- Desktop layout uses a left rule column plus a right workspace.
- The entire left column is dedicated to rules. Rule actions such as global preview and add may live in the rule-column header; the rest of the column is rule cards only.
- If rules exceed available height, the left rule column scrolls internally.
- Rule cards include a moderately sized product thumbnail on the left and rule summary on the right. The image and text stack must share the same card height.
- Rule-list thumbnails are smaller than search-preview product images but large enough for visual recognition.
- The right workspace top row contains active rule editing and compact Telegram delivery summary.
- The active rule editor takes about 70% of the top row and the Telegram summary takes about 30%.
- Active rule fields use two rows: keyword and filter words first; price limit, comment threshold, worthy-rate threshold, and schedule interval second.
- Filter words use an input plus compact removable chips. The chip area has fixed maximum height, shows up to two rows, and scrolls internally when more filters exist.
- The active rule editor has a local bottom action area: Preview/Search is the secondary validation action, Save is the primary commit action.
- The search preview spans the full right workspace width below both the active rule editor and Telegram summary.
- Desktop search preview uses 3 columns x 2 rows of independent wide result cards. Each product remains self-contained with its own image, title, metrics, and open-link action.
- Search result cards are compact wide validation cards, not large product-detail cards: fixed left media, right-side title, stable price/comment/worthy-rate metrics, and clear open-link action.
- Product images use fixed dimensions and stable placeholders when loading fails: rule cards show the keyword initial; search results show the product-title initial.
- Image placeholders must match loaded image dimensions, avoid broken-image icons, avoid grayscale/blurred real images, and prevent layout shift.
- Visual design uses one coherent theme, one restrained accent color, consistent radius, visible labels, 44px minimum interactive targets, and tabular numbers for price/comment/worthy-rate metrics.
- Status and action colors are limited to three semantic meanings: green for enabled/healthy delivery, one primary accent for save/primary flow, and red for delete/error. Other controls remain neutral.
- Destructive actions use danger styling and must be visually separated from primary commit actions.
- Avoid over-colorful per-icon styling, decorative status dots, generic three-equal-card layout, excessive glows/gradients, placeholder-only fields, disconnected floating buttons, and large saturated status blocks.
- Mobile order is rules, active rule editor, search preview, then Telegram summary. Telegram is last on mobile because notification settings are mostly fixed after setup.
- Mobile must avoid horizontal page overflow; rule cards and the result preview may use contained horizontal scrolling or compact stacking within their own sections.
- Dark mode must continue to follow system preference.

## Acceptance Criteria

- [ ] Desktop screenshot shows a thin top bar, a left rule column, a right workspace top row with active rule editor and compact Telegram summary, and a 3 x 2 independent wide-card search preview spanning the workspace below.
- [ ] Desktop screenshot shows rule creation/editing as the primary hierarchy, temporary search preview as validation, and Telegram as secondary delivery readiness.
- [ ] Rule card screenshot shows left thumbnail and right rule summary aligned to the same height.
- [ ] Filter chip stress case with more than two rows scrolls inside the chip region without pushing neighboring fields or changing the editor height.
- [ ] Search preview shows six independent compact wide result cards on desktop and does not render as large product-detail cards.
- [ ] Image load failure renders stable initial placeholders without layout shift.
- [ ] Mobile browser measurement confirms `scrollWidth <= clientWidth` and content order is rules, editor, preview, Telegram.
- [ ] Icon-only buttons have `aria-label` and `title`.
- [ ] Inline script passes `node --check`.
- [ ] `go test ./...` passes.
- [ ] `git diff --check` has no whitespace errors except acceptable Windows line-ending warnings.
- [ ] `data/users.db` is not committed.

## Out Of Scope

- Backend API changes.
- Restoring user data.
- Adding providers beyond Telegram.
- Deploying before local screenshot verification.
