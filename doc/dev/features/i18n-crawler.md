# Notes related to ./i18n/crawler

## Primer on html and jet's treatment of templates

```
<a>abc {{.Hi}}</a>
```

html: `abc {{.Hi}}` is text

jet: `<a>abc ` is text
jet: `</a>` is text

So, jet doesn't care about HTML.

And we have to care about `{* *}` comments.

## Coalesce inline tags

Example: in below, the `<a>` tag should be part of the string.

```html
Log in with your Pixiv account's cookie to access features above. To learn how to obtain your cookie, please see <a href="https://pixivfe-docs.pages.dev/obtaining-pixivfe-token/">the guide on obtaining your PixivFE token</a>.
```

Tags to consider inline:

- `a`