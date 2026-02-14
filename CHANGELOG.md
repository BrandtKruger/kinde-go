<!-- Ideally, this should get auto-generated via tools like [auto-changelog](https://github.com/CookPete/auto-changelog). Eventually, this will get set up as part of the repository template. -->

## Unreleased

### Breaking changes

- **Management API (generated):** Response types for the delete API application scope operation were renamed to fix a typo: `DeleteAPIAppliationScope*` → `DeleteAPIApplicationScope*` (e.g. `DeleteAPIAppliationScopeBadRequest` → `DeleteAPIApplicationScopeBadRequest`). If you reference the old type names in your code, update them to the new spelling.
