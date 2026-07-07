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
- Latest screenshot feedback, 2026-07-07: current cards still waste vertical space, especially below the editor and inside result cards; upper text hierarchy is too small; rule search/delete actions should sit on the right of the rule-card header; comment and worthy metrics should use SMZDM-like inline icons; metric outer boxes, especially the comment box, are visually rejected; black circular field icons in the light editor/TG form are too abrupt.
- Latest layout feedback, 2026-07-07: the lower empty area in the rule editor is unnecessary and should be removed. Telegram is visually too small, but its controls are poorly distributed vertically; link preview can be made compact and placed on the same row as Parse Mode / HTML instead of occupying a full separate block.
- Latest title/status feedback, 2026-07-07: page title should be "什么值得买商品提醒". Rule title and enabled indicator can be larger. A save action on the far right of the rule-detail title bar is not useful. Telegram enable state should be represented by a Telegram plane icon: enabled uses a filled blue background; disabled uses a blank/empty state.
- Latest composition feedback, 2026-07-07: the dashboard may remain divided into functional blocks, but it must read as one integrated product surface. Avoid obvious heavy split lines, black gutters, and isolated floating islands that make each block look unrelated.
- Latest icon feedback, 2026-07-07: icon colors are too scattered. Use a disciplined semantic icon color system instead of assigning a different color to each icon.
- Latest feedback, 2026-07-08: top-bar Telegram uses the official logo and no text; Telegram panel itself uses a real switch control. Field icons should not use black square/circle backgrounds. Left rules must scroll internally when there are more than five rules. Search results should use space better, show up to eight visible result cards, allow product titles to wrap, and remove per-card open icons because the title is already clickable; keep only the preview panel's global open-web button.

## Requirements

- Preserve existing DOM IDs, JavaScript behavior, route contracts, and feature coverage unless a later implementation note proves a wrapper change is required.
- Desktop uses a thin global top bar, about 56-64px, for product name, save status, Telegram health, and theme state only. It must not become a hero/header block.
- The product/page title should read "什么值得买商品提醒".
- Desktop layout uses a left rule column plus a right workspace.
- Desktop layout should use one integrated app canvas. Functional blocks may remain, but their separation must come from soft spacing, subtle surface changes, and shared alignment rather than strong divider lines or black gutters.
- Use a shared app canvas with soft internal panels. Avoid distinct card-like islands separated by heavy visible gaps.
- The entire left column is dedicated to rules. Rule actions such as global preview and add may live in the rule-column header; the rest of the column is rule cards only.
- If rules exceed available height, the left rule column scrolls internally.
- Rule cards include a moderately sized product thumbnail on the left and rule summary on the right. The image and text stack must share the same card height.
- Rule-list thumbnails are smaller than search-preview product images but large enough for visual recognition.
- The right workspace top row contains active rule editing and compact Telegram delivery summary.
- The active rule editor takes about 70% of the top row and the Telegram summary takes about 30%.
- Active rule fields use two rows: keyword and filter words first; price limit, comment threshold, worthy-rate threshold, and schedule interval second.
- Keep the current rule-detail field layout rhythm instead of forcing a pure 2-column or pure 3-column grid. Keyword and filter fields remain wider; numeric/threshold/runtime fields keep the current compact grid relationship.
- Scan interval and per-push count stay visible in the rule-detail form. They are not folded into a hidden advanced section.
- Filter words use an input plus compact removable chips. The chip area has fixed maximum height, shows up to two rows, and scrolls internally when more filters exist.
- The active rule editor has a local bottom action area: Preview/Search is the secondary validation action, Save is the primary commit action.
- Rule-detail actions belong to the rule-detail form footer. The footer contains a compact action group: secondary "preview current rule" and primary "save and enable". Do not split these actions left/right across the panel, and do not place the trigger in the search-preview section.
- The active rule editor should shrink to actual content height. It must not reserve a large empty bottom region just to fill the panel.
- The search preview spans the full right workspace width below both the active rule editor and Telegram summary.
- Search preview should size to actual result count instead of reserving empty slots. Show up to eight results in the first visible area; when more than eight results exist, the preview list scrolls internally.
- Desktop search preview uses independent product cards. Each product remains self-contained with its own image, wrapping title, and metrics. The product title is clickable; do not render a separate open-link action inside every card.
- Search result cards are compact wide validation cards, not large product-detail cards: fixed left media, right-side title, stable price/comment/worthy-rate metrics, and clear open-link action.
- Search result typography must make the product title and price easy to read at a glance; do not reserve large blank lower areas inside cards.
- Search-result open-web action appears only once in the search-preview title bar as the global open-web icon button.
- Rule card actions for preview/search and delete align to the right side of the rule title row.
- Rule enabled state should render as an icon in the rule list instead of text such as "已启用".
- Rule titles and enabled-state icons should be larger than the current small status text treatment.
- Rule-card metrics must use inline metadata, not large bordered metric boxes. Comment and worthy-rate should read as small icon+number metadata in a style close to SMZDM's own product metadata.
- Search-result metrics should use the same inline metadata language as rule cards, with price visually strongest and comment / worthy-rate secondary.
- Form field icons should be preserved, but they must be the same visual size as the field label text, not large standalone circular badges. Icons must be semantically appropriate for each field and use theme-aware background/foreground colors for light and dark mode.
- Telegram configuration should be more legible than the current narrow/tiny card, while keeping secondary priority. Its layout should be compact: Bot Token and Chat ID remain primary fields; Parse Mode / HTML and link preview should share a compact row when space allows.
- Telegram's main enable switch belongs in the panel title bar. The body should contain only four configuration controls: Bot Token, Chat ID, Parse Mode / HTML, and link preview.
- Top-bar Telegram control should use the official Telegram logo. Enabled state uses normal color / filled emphasis; disabled state is greyed. The Telegram panel enable control itself uses a normal switch.
- The rule detail/editor panel should not contain a rule enable switch. Saving a rule automatically enables it.
- Product images use fixed dimensions and stable placeholders when loading fails: rule cards show the keyword initial; search results show the product-title initial.
- Image placeholders must match loaded image dimensions, avoid broken-image icons, avoid grayscale/blurred real images, and prevent layout shift.
- Visual design uses one coherent theme, one restrained accent color, consistent radius, visible labels, 44px minimum interactive targets, and tabular numbers for price/comment/worthy-rate metrics.
- Panel boundaries should be quiet and coherent: shared background, consistent radius, low-contrast borders, and aligned gutters. Avoid heavy vertical/horizontal split lines that visually break the interface into unrelated cards.
- Status and action colors are limited to three semantic meanings: green for enabled/healthy delivery, one primary accent for save/primary flow, and red for delete/error. Other controls remain neutral.
- Icon colors must follow the same semantic system: primary blue for primary/search/open actions, green for enabled/healthy state, red for delete/error, and neutral muted color for field icons/secondary metadata. Avoid per-icon random colors.
- Destructive actions use danger styling and must be visually separated from primary commit actions.
- Avoid over-colorful per-icon styling, decorative status dots, generic three-equal-card layout, excessive glows/gradients, placeholder-only fields, disconnected floating buttons, and large saturated status blocks.
- Mobile order is rules, active rule editor, search preview, then Telegram summary. Telegram is last on mobile because notification settings are mostly fixed after setup.
- Mobile must avoid horizontal page overflow; rule cards and the result preview may use contained horizontal scrolling or compact stacking within their own sections.
- Dark mode must continue to follow system preference.

## Acceptance Criteria

- [ ] Desktop screenshot shows a thin top bar, a left rule column, a right workspace top row with active rule editor and compact Telegram summary, and an independent-card search preview spanning the workspace below.
- [ ] Desktop screenshot reads as one integrated app canvas, not separate islands divided by heavy black gutters or strong split lines.
- [ ] Page title reads "什么值得买商品提醒".
- [ ] Desktop screenshot shows rule creation/editing as the primary hierarchy, temporary search preview as validation, and Telegram as secondary delivery readiness.
- [ ] Rule card screenshot shows left thumbnail and right rule summary aligned to the same height.
- [ ] Filter chip stress case with more than two rows scrolls inside the chip region without pushing neighboring fields or changing the editor height.
- [ ] Search preview sizes to actual result count, shows up to eight independent compact product cards before internal scrolling, and does not render as large product-detail cards.
- [ ] Search result cards do not show large unused blank areas under the main content, and product title/price hierarchy is visibly stronger than secondary metadata.
- [ ] Search-result cards do not render per-card open-web icons; only the search-preview header has the global open-web icon.
- [ ] Rule cards place preview/search and delete actions on the right side of the title row.
- [ ] Rule enabled state renders as an icon in the rule list, not as the text "已启用".
- [ ] Rule title and enabled-state icon are visually larger than the previous small text status.
- [ ] Rule-card and search-result comment / worthy-rate metrics render as inline icon+number metadata without heavy outer metric boxes.
- [ ] Rule editor and Telegram form preserve field icons, but icons are label-sized, semantically appropriate, aligned with label text, and theme-aware in light/dark mode.
- [ ] Rule editor contains no large unused bottom blank region; its actions sit close to the actual form content.
- [ ] Rule-detail footer contains the current-rule preview action and the primary save-and-enable action together; the search-preview section contains results only.
- [ ] Telegram panel is legible and compact: Parse Mode / HTML and link preview sit in one row on desktop.
- [ ] Telegram panel enable control is a real switch in the Telegram panel title bar, and the panel body contains only Bot Token, Chat ID, Parse Mode / HTML, and link preview.
- [ ] Top-bar Telegram control uses the official Telegram logo; enabled appears normal/color-emphasized and disabled appears greyed.
- [ ] Icon colors are visually consistent: primary/search/open actions, enabled/healthy state, delete/error, and neutral metadata are the only color categories.
- [ ] Rule detail/editor contains no enable switch; saving a rule automatically marks it enabled.
- [ ] Scan interval and per-push count remain directly visible and editable in the rule-detail form.
- [ ] Rule-detail field layout keeps the current mixed grid rhythm: wider keyword/filter fields and compact aligned threshold/runtime fields.
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

## Decisions

- Metrics presentation: use inline SMZDM-like metadata for price, worthy-rate, and comment counts. Do not use large bordered metric boxes for rule cards or search result cards.
- Telegram layout: put the main enable switch in the title bar. Keep only four body controls: Bot Token, Chat ID, Parse Mode / HTML, and link preview.
- Field icon treatment: keep icons in the editor and Telegram forms, but make them label-sized, semantically matched to the field, and theme-aware instead of black circular badges.
- Search preview density: size the preview to actual result count; show up to eight independent product cards before internal scrolling.
- Rule enabled behavior: show enabled state as an icon in rule-list cards. Remove the enable switch from rule details; saving a rule automatically enables it.
- Product title: use "什么值得买商品提醒".
- Rule detail actions: do not put save on the far right of the rule-detail title bar.
- Rule detail action placement: put preview-current-rule and save-and-enable together in the rule-detail form footer. Search preview is display-only.
- Telegram enabled visual: use Telegram plane icon with filled blue enabled state and blank disabled state.
- Rule-detail low-frequency fields: keep scan interval and per-push count visible in the main form.
- Rule-detail field grid: preserve the current mixed layout rhythm rather than changing to a strict 2-column or 3-column grid.
- Search-result action placement: do not put open-web actions on every card. Keep only the global open-web action in the search-preview title bar.
- Overall composition: keep functional blocks, but make the dashboard read as one integrated canvas using soft gutters and low-contrast boundaries instead of heavy separators.
- Icon color system: use semantic categories only. Primary/search/open = blue, enabled/healthy = green, destructive/error = red, fields/secondary metadata = neutral.

## Open Questions

- None.
