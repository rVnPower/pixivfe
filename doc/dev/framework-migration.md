# Migrating from gofiber to net/http -- the plan

## Templating

First, convert all `c.Render("about", ...)` to `Render(c, Data_about{...}`.
Then, the templating engine is decoupled from gofiber.

### Fixing current templates

gofiber doesn't mention this crucial distinction of Jet: the difference between variables (e.g. PageURL) and data (the data parameter to `Render(...)`).

Data access must be prefixed with a dot. `.Illust` is valid. `Illust` is a variable but not data.

Current templates use the variable style (`Illust`), but that is wrong.
