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

## Tips

To access `/:abc`, use `r.PathValue("abc")`.

Idea: we create a compat layer of {w, r} that has the same API as *fiber.Ctx.
