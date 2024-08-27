# Migrating from gofiber to net/http -- the plan

## Tips

- To access `/:abc`, use `r.PathValue("abc")`.
- Want to associate arbitrary value with a request? Use `r.Context.Value()`.
- Don't use `r.Response`. It's `nil`.
- Only use the following HTTP status codes:
  - 303 StatusSeeOther: set method to GET
  - 307 StatusTemporaryRedirect: method and body not changed
  - 308 StatusPermanentRedirect: method and body not changed

## todo

- add limiter (maybe it should be in nginx)
- add caching (maybe it should be in nginx)

## Guide for caching

The majority of the site's traffic is images, and all images on the Pixiv route can be cached forever.
  Cached in browser: pixiv headers cache images for year
  On the server, the reverse proxy (e.g. nginx) can be cached worry-free forever.

JSON requests to Pixiv can be cached for some time.
  todo: We should have a middleware handling that. 

Assets are already cached properly thanks to net/http.

Every rendered page can be cached for 5 minutes or so in the browser.
However, after changing settings, the page need to be refreshed instantly.
  Cache in browser: todo: read below implement in template/render.go
  The reverse proxy shouldn't do anything about this.

The only three headers that controls caching in browser and middleboxes:

- [Age](https://httpwg.org/specs/rfc9111.html#field.age)
- [Cache-Control](https://httpwg.org/specs/rfc9111.html#field.cache-control)
- [Expires](https://httpwg.org/specs/rfc9111.html#field.expires)
