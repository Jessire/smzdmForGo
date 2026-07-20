# Design

- Use an anchored popover positioned from `#addRule`, with the existing three choices.
- Keep shared filter controls in the main rule editor and hide only the product keyword field for system rules.
- Persist discovery shared filters in `GlobalHotConfig`, preserving old JSON values through zero-value defaults.
- Apply shared filter matching after type-specific hot/author matching.
- Include `referral` in preview cards and keep an explicit empty-state explanation when no item passes the active threshold.
