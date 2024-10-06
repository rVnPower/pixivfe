Primer on html and jet's treatment of templates

```
<a>abc {{.Hi}}</a>
```

html: `abc {{.Hi}}` is text

jet: `<a>abc ` is text
jet: `</a>` is text

So, jet doesn't care about HTML.

And we have to care about `{* *}` comments.
