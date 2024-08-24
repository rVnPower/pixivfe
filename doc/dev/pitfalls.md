## Strict CSP

Reference: search for "Content-Security-Policy" in **.go

Current CSP disallows inline styles and scripts and iframes.

## Low Quality Go Module: net/url

`url.Path` is stored decoded (no %XX). `url.Scheme` is stored without `://` (mandated by RFC). Not sure why Go does that. Felt like this is bound to cause some nasty bug on decoding and encoding.

Current proxied URLs don't have weird characters in them. Hopefully it stays this way.

Solution: Replace "net/url" with a better third-party module

## Jet templating engine is not type checked

Solution: [templ](https://github.com/a-h/templ)

## Jet `.Illust` vs `Illust`

Be careful when you write a .jet.html template file.

If you see error like `<variable name> not found in map[...]` when visiting a page, you need to add a dot before all data member access. e.g. from `Illust` to `.Illust`. gofiber doesn't mention this.

In Jet, variables are accessed without a dot (`PageURL`), while data members (the `data` parameter to `Render(...)`). are accessed with a dot in front.

There are 3 types of variables.

- Global variables. Global functions are those.
- Page variables. They are provided to Jet by the server.
- Temporary variables. You use `abc := ...` to define those.
