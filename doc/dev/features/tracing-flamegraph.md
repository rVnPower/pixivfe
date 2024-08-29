# About tracing in PixivFE

Every request to pixiv websites should go through core/requests.go.

Every request to pixiv websites is traced.
Every server request is traced.

## How to see flamegraph

Run PixivFE **in dev mode**, visit some pages, then visit URL /diagnostics.

## Useful for fixing Vega-Lite

https://vega.github.io/vega-lite/examples/interactive_legend.html
https://vega.github.io/editor/#/examples/vega-lite/interactive_legend
https://vega.github.io/vega-lite/docs/tooltip.html
