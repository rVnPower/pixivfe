# Migrating from gofiber to net/http -- the plan

- Config [already decoupled]
- Templating [decoupled, waiting for integration]
- Logging
- Router 
  features
    - /users/:id/:category? (optional path segment)
    - /i.pximg.net/* (wildcard)
- Middleware
  - Rate limit (optional, could be loosely-coupled)
  - Caching (optional, could be loosely-coupled)
