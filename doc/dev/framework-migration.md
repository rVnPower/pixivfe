# Migrating from gofiber to net/http -- the plan

## Templating

First, convert all `c.Render("about", ...)` to `Render(c, Data_about{...}`.
Then, the templating engine is decoupled from gofiber.

### Fixing current templates

If you see error like `<variable name> not found in map[...]` when visiting a page, you need to add a dot before all data member access. e.g. from `Illust` to `.Illust`.

gofiber doesn't mention this.
In Jet, variables are accessed without a dot (`PageURL`), while data members (the `data` parameter to `Render(...)`). are accessed with a dot in front.

Current templates use the variable style (`Illust`), but that is wrong.
