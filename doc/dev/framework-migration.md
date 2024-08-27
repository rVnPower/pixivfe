# Migrating from gofiber to net/http -- the plan

## Tips

- To access `/:abc`, use `r.PathValue("abc")`.
- Want to associate arbitrary value with a request? Use `r.Context.Value()`.
- Don't use `r.Response`. It's `nil`.

## todo

- correct redirect status codes. currently i put in whatever.
  - 303 StatusSeeOther: set method to GET
  - 307 Temporary: method and body not changed
  - 308 Permanent: method and body not changed
- add limiter (maybe it should be in nginx)
- add caching (maybe it should be in nginx)
