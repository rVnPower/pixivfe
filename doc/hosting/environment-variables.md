# Environment Variables

PixivFE's behavior is controlled by environment variables. Currently, you can only set variables directly in your environment.

An example configuration is provided in [`.env.example`](https://codeberg.org/VnPower/PixivFE/src/branch/v2/.env.example).

!!! tip
    To quickly set up PixivFE, you need to define two required environment variables:

    - `PIXIVFE_TOKEN`: Your Pixiv account cookie, which is necessary for accessing Pixiv's Ajax API. Refer to the [guide on obtaining the PIXIVFE_TOKEN cookie](obtaining-pixivfe-token.md) for details on how to acquire your Pixiv token.
    - `PIXIVFE_PORT`: The port number on which PixivFE will run, for example, `8282`.

    For basic usage, configure your environment variables as follows:
    ```
    PIXIVFE_TOKEN=123456_AaBbccDDeeFFggHHIiJjkkllmMnnooPP
    PIXIVFE_PORT=8282
    ```

    If you are setting up a development environment, enable the development mode by also setting:
    ```
    PIXIVFE_DEV=true
    ```

## Required variables

### `PIXIVFE_PORT` or `PIXIVFE_UNIXSOCKET`

**Required**: Yes (one of the two)

- `PIXIVFE_PORT`: Port to listen on, e.g., `PIXIVFE_PORT=8282`.
- `PIXIVFE_UNIXSOCKET`: [UNIX socket](https://en.wikipedia.org/wiki/Unix_domain_socket) to listen on, e.g., `PIXIVFE_UNIXSOCKET=/srv/http/pages/pixivfe`.

### `PIXIVFE_TOKEN`

**Required**: Yes

Your Pixiv account cookie, used by PixivFE for authorization to fully access Pixiv's Ajax API.

See the [Obtaining the `PIXIVFE_TOKEN` cookie](obtaining-pixivfe-token.md) guide for detailed instructions.

## Optional variables

### `PIXIVFE_HOST`

**Required**: No (ignored if `PIXIVFE_UNIXSOCKET` is set)

!!!note
    If you're **not using a reverse proxy** or **running PixivFE inside Docker**, you should set `PIXIVFE_HOST=0.0.0.0`. This will allow PixivFE to accept connections from any IP address or hostname. If you don't set this, PixivFE will refuse direct connections from other machines or devices on your network.

This setting specifies the hostname or IP address that PixivFE should listen on and accept incoming connections from. For example, if you want PixivFE to only accept connections from the same machine (your local computer), you can set `PIXIVFE_HOST=localhost`.

### `PIXIVFE_REQUESTLIMIT`

**Required**: No

Set to a number to enable the built-in rate limiter, e.g., `PIXIVFE_REQUESTLIMIT=15`.

It's recommended to enable rate limiting in the reverse proxy in front of PixivFE rather than using this.

### `PIXIVFE_IMAGEPROXY`

**Required**: No, defaults to using the built-in proxy

!!! note
    The protocol **must** be included in the URL, e.g., `https://piximg.example.com`, where `https://` is the protocol used.

The URL of the image proxy server. Pixiv requires `Referer: https://www.pixiv.net/` in the HTTP request headers to fetch images directly. Set this variable if you wish to use an external image proxy or are unable to get images directly from Pixiv.

See [hosting an image proxy server](image-proxy-server.md) or the [list of public image proxies](../public-image-proxies.md).

### `PIXIVFE_USERAGENT`

**Required**: No

**Default:** `Mozilla/5.0 (Windows NT 10.0; rv:123.0) Gecko/20100101 Firefox/123.0`

The value of the `User-Agent` header used for requests to Pixiv's API.

### `PIXIVFE_ACCEPTLANGUAGE`

**Required**: No

**Default:** `en-US,en;q=0.5`

The value of the `Accept-Language` header used for requests to Pixiv's API. Change this to modify the response language.

### `PIXIVFE_PROXY_CHECK_INTERVAL`

**Required**: No

**Default:** `8h`

The interval in minutes between proxy checks. Defaults to 8 hours if not set.
Please specify this value in Go's `time.Duration` notation, e.g. `2h3m5s`.
You can disable this by setting the value to 0. Then, proxies will only be checked once at server initialization.

### `PIXIVFE_TOKEN_LOAD_BALANCING`

**Required**: No

**Default:** `round-robin`

Specifies the method for selecting tokens when multiple tokens are provided in `PIXIVFE_TOKEN`.

Valid options:

- `round-robin`: Tokens are used in a circular order.
- `random`: A random token is selected for each request.

This option is useful when you have multiple Pixiv accounts and want to distribute the load across them.

### `PIXIVFE_DEV`

**Required**: No

Set to any value to enable development mode, in which the server will live-reload HTML templates + SCSS files and disable caching, e.g., `PIXIVFE_DEV=true`.
