# Frontend Spec

## Scope

Applies to `template/html/index.html`, which contains the Web panel HTML, CSS, and inline JavaScript.

## Pre-Development Checklist

- Read `.trellis/spec/guides/index.md`.
- Preserve the four user-facing areas: 商品规则, 搜索结果, 规则详情, Telegram 通知.
- Verify desktop and mobile layout with a real browser after CSS changes.
- Keep touch targets near 44px or larger for icon-only actions.

## Local Patterns

- The current panel uses the `body.zero-redesign` class and the `Zero redesign v3` CSS block as the active design layer.
- Desktop information flow should keep rule selection and search results close. The preferred desktop grid is `rules preview editor`, with Telegram notification configuration under rule details.
- Mobile layout stacks as `rules`, `editor`, then the inspector column containing preview and notification.
- Icon-only controls must include `aria-label` and `title`; existing buttons such as `#searchAllRules`, `#addRule`, `#searchRule`, and `#saveProductConfig` show this pattern.
- Product images are rendered through search results and image proxy behavior. CSS must not apply grayscale or blur filters to product thumbnails.

## Avoid

- Do not place low-frequency notification settings in the primary center column on desktop.
- Do not force all panels to the same viewport height when their content is shorter.
- Do not make the user jump from the left rule list to the far-right column to inspect search results.
- Do not introduce wxPusher, DingTalk, or sign-in UI; current notification UI is Telegram-only.
