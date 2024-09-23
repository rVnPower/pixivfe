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

## `PIXIVFE_PORT` or `PIXIVFE_UNIXSOCKET`

**Required**: Yes (one of the two)

- `PIXIVFE_PORT`: Port to listen on, e.g., `PIXIVFE_PORT=8282`.
- `PIXIVFE_UNIXSOCKET`: [UNIX socket](https://en.wikipedia.org/wiki/Unix_domain_socket) to listen on, e.g., `PIXIVFE_UNIXSOCKET=/srv/http/pages/pixivfe`.

## `PIXIVFE_TOKEN`

**Required**: Yes

Your Pixiv account cookie, used by PixivFE for authorization to fully access Pixiv's Ajax API. This variable can contain multiple tokens separated by commas, which is useful for load balancing across multiple Pixiv accounts.

Example:
```
PIXIVFE_TOKEN=123456_AaBbccDDeeFFggHHIiJjkkllmMnnooPP,789012_QqRrSsTtUuVvWwXxYyZz
```

See the [Obtaining the `PIXIVFE_TOKEN` cookie](obtaining-pixivfe-token.md) guide for detailed instructions.

## `PIXIVFE_HOST`

**Required**: No (ignored if `PIXIVFE_UNIXSOCKET` is set)

!!!note
    If you're **not using a reverse proxy** or **running PixivFE inside Docker**, you should set `PIXIVFE_HOST=0.0.0.0`. This will allow PixivFE to accept connections from any IP address or hostname. If you don't set this, PixivFE will refuse direct connections from other machines or devices on your network.

This setting specifies the hostname or IP address that PixivFE should listen on and accept incoming connections from. For example, if you want PixivFE to only accept connections from the same machine (your local computer), you can set `PIXIVFE_HOST=localhost`.

## `PIXIVFE_REPO_URL`

**Required**: No

**Default**: `https://codeberg.org/VnPower/PixivFE`

The URL of the PixivFE source code repository. This is used in the about page to provide links to the project's source code and specific commit information. You can change this if you're running a fork of PixivFE and want to link to your own repository instead.

## `PIXIVFE_REQUESTLIMIT`

**Required**: No

Set to a number to enable the built-in rate limiter, e.g., `PIXIVFE_REQUESTLIMIT=15`.

It's recommended to enable rate limiting in the reverse proxy in front of PixivFE rather than using this.

## `PIXIVFE_IMAGEPROXY`

**Required**: No, defaults to using the built-in proxy

!!! note
    The protocol **must** be included in the URL, e.g., `https://piximg.example.com`, where `https://` is the protocol used.

The URL of the image proxy server. Pixiv requires `Referer: https://www.pixiv.net/` in the HTTP request headers to fetch images directly. Set this variable if you wish to use an external image proxy or are unable to get images directly from Pixiv.

See [hosting an image proxy server](image-proxy-server.md) or the [list of public image proxies](../public-image-proxies.md).

## `PIXIVFE_USERAGENT`

**Required**: No

**Default:** `Mozilla/5.0 (Windows NT 10.0; rv:123.0) Gecko/20100101 Firefox/123.0`

The value of the `User-Agent` header used for requests to Pixiv's API.

## `PIXIVFE_ACCEPTLANGUAGE`

**Required**: No

**Default:** `en-US,en;q=0.5`

The value of the `Accept-Language` header used for requests to Pixiv's API. Change this to modify the response language.

## `PIXIVFE_PROXY_CHECK_INTERVAL`

**Required**: No

**Default:** `8h`

The interval in minutes between proxy checks. Defaults to 8 hours if not set.
Please specify this value in Go's `time.Duration` notation, e.g. `2h3m5s`.
You can disable this by setting the value to 0. Then, proxies will only be checked once at server initialization.

## `PIXIVFE_TOKEN_LOAD_BALANCING`

**Required**: No

**Default:** `round-robin`

Specifies the method for selecting tokens when multiple tokens are provided in `PIXIVFE_TOKEN`.

Valid options:

- `round-robin`: Tokens are used in a circular order.
- `random`: A random token is selected for each request.
- `least-recently-used`: The token that hasn't been used for the longest time is selected.

This option is useful when you have multiple Pixiv accounts and want to distribute the load across them, reducing the risk of rate limiting for individual accounts by the Pixiv API.

## Exponential backoff configuration

PixivFE implements exponential backoff for API requests and token management to handle failures gracefully and manage rate limiting. The following environment variables can be used to configure this behavior, fine-tuning the exponential backoff behavior for both API requests and token management. If not set, the default values will be used.

For more detailed information about the implementation of exponential backoff in PixivFE, please refer to the [Exponential Backoff documentation](../dev/features/exponential_backoff.md).

### API request level backoff

These settings control how PixivFE handles retries for individual API requests. The backoff time starts at the base timeout and doubles with each retry, up to the maximum backoff time.

#### `PIXIVFE_API_MAX_RETRIES`

**Required**: No

**Default:** `3`

Maximum number of retries for API requests.

#### `PIXIVFE_API_BASE_TIMEOUT`

**Required**: No

**Default:** `500ms`

Base timeout duration for API requests.

#### `PIXIVFE_API_MAX_BACKOFF_TIME`

**Required**: No

**Default:** `8000ms`

Maximum backoff time for API requests.

### Token management level backoff

These settings control how PixivFE manages token timeouts when a token encounters repeated failures. The backoff time for a token starts at the base timeout and doubles with each failure, up to the maximum backoff time.

#### `PIXIVFE_MAX_RETRIES`

**Required**: No

**Default:** `5`

Maximum number of retries for token management.

#### `PIXIVFE_BASE_TIMEOUT`

**Required**: No

**Default:** `1000ms`

Base timeout duration for token management.

#### `PIXIVFE_MAX_BACKOFF_TIME`

**Required**: No

**Default:** `32000ms`

Maximum backoff time for token management.

## Network proxy configuration

Used to set the [proxy server](https://en.wikipedia.org/wiki/Proxy_server) that PixivFE will use for all requests. Not to be confused with the image proxy, which is used to comply with the `Referer` check required by `i.pximg.net`.

Requests use the proxy specified in the environment variable that matches the scheme of the request (`HTTP_PROXY` or `HTTPS_PROXY`). This selection is based on the scheme of the **request being made**, not on the protocol used by the proxy server itself.

### `HTTPS_PROXY`

**Required**: No

Proxy server used for requests made over HTTPS.

### `HTTP_PROXY`

**Required**: No

Proxy server used for requests made over plain HTTP.

## Development options

### `PIXIVFE_DEV`

**Required**: No

Set to any value to enable development mode, e.g., `PIXIVFE_DEV=true`. In development mode:

1. The server will live-reload HTML templates and SCSS files.
2. Caching is disabled.
3. Additional debug information is logged.
4. Responses are saved to `PIXIVFE_RESPONSE_SAVE_LOCATION`.

This setting is useful for developers working on PixivFE itself or for troubleshooting issues in a development environment.

### `PIXIVFE_RESPONSE_SAVE_LOCATION`

**Required**: No

**Default**: `/tmp/pixivfe/responses`

Defines where responses from the Pixiv API are saved when in development mode.
