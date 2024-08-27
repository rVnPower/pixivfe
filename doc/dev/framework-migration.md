# Migrating from gofiber to net/http -- the plan

- Config [already decoupled]
- Templating [decoupled, waiting for integration]
- Router 
  features
    - /users/:id/:category? (optional path segment)
    - /i.pximg.net/* (wildcard)
- Middleware
  - Logging
  - Rate limit (optional, could be loosely-coupled)
  - Caching (optional, could be loosely-coupled)

## Problem

net/http handlers don't return errors. We have to make our own ServeMux that allows functions to return `error`, possibly.

net/http expects handlers to panic on error, while we don't panic. We need to log the errors anyway.

Idea: we create a compat layer of {w, r} that has the same API as *fiber.Ctx.

## Tips

- To access `/:abc`, use `r.PathValue("abc")`.
- Want to associate arbitrary value with a request? Use `r.Context.Value()`.
- Don't use `r.Response`. It's `nil`.

## todo

- correct redirect status codes. currently i put in whatever.
