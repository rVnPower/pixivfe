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


## Tips

To access `/:abc`, use `r.PathValue("abc")`.
