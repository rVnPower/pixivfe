# Migrating from gofiber to net/http -- the plan

Templating
  First, convert all `c.Render("about", ...)` to `Render(c, Data_about{...}`.
  Then, the templating engine is decoupled from gofiber.
