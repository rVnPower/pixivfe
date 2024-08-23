# PixivFE

A privacy-respecting alternative front-end for Pixiv that doesn't suck.

<p>
<a href="https://codeberg.org/vnpower/pixivfe">
<img alt="Get it on Codeberg" src="https://get-it-on.codeberg.org/get-it-on-blue-on-white.png" height="60">
</a>
</p>

![CI badge](https://ci.codeberg.org/api/badges/12556/status.svg)
[![Go Report Card](https://goreportcard.com/badge/codeberg.org/vnpower/pixivfe/v2)](https://goreportcard.com/report/codeberg.org/vnpower/pixivfe)

Questions? Feedback? You can [PM me](https://matrix.to/#/@vnpower:eientei.org) on Matrix! You can also see the [Known quirks](https://pixivfe.pages.dev/known-quirks/) page to check if your issue has a known solution.

You can keep track of this project's development using the [roadmap](doc/dev/general.md).

## Features

- Lightweight - both the interface and the code
- Privacy-first - the server will do the work for you
- No bloat - we only serve HTML, CSS and minimal JS code
- Open source - you can trust me!

## Hosting

You can use PixivFE for personal use! Assuming that you use an operating system that can run POSIX shell scripts, install `go`, clone this repository, modify the `run.sh` file, and profit!
I recommend self-hosting your own instance for personal use, instead of relying entirely on official instances.

To deploy PixivFE using Docker or the compiled binary, see [Hosting PixivFE](https://pixivfe.pages.dev/hosting-pixivfe/).

PixivFE can work with or without an external image proxy server. Here is the [list of public image proxies](https://pixivfe.pages.dev/public-image-proxies/).
See [hosting a Pixiv image proxy](https://pixivfe.pages.dev/hosting-image-proxy-server/) if you want to host one yourself.


## Development

**Requirements:**

- [Go](https://go.dev/doc/install) (to build PixivFE from source)
- [Sass](https://github.com/sass/dart-sass/) (will be run by PixivFE in development mode)

To install Dart Sass, you can choose any of the following methods.

- use system package manager (usually called `dart-sass`)
- download executable from [the official release page](https://github.com/sass/dart-sass/releases)
- `pnpm i -g sass`

```bash
# Clone the PixivFE repository
git clone https://codeberg.org/VnPower/PixivFE.git && cd PixivFE

# Run in PixivFE in development mode (styles and templates reload automatically)
PIXIVFE_DEV=1 <other_environment_variables> go run .
```

## Instances

<!-- The current instance table is really wide; maybe there's a better way of formatting it without losing information?
The badges are also difficult to read on a small screen due to Codeberg shrinking the width of the columns -->

| Name              | URL                                             | Country | Cloudflare? | Observatory Grade                                                                                                                               | Status                                                                                                                                               |
|-------------------|-------------------------------------------------|---------|-------------|-------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------|
| exozyme (Official)| [https://pixivfe.exozy.me](https://pixivfe.exozy.me) | US      | No          | [![Mozilla HTTP Observatory Grade](https://img.shields.io/mozilla-observatory/grade-score/pixivfe.exozy.me?label=)](https://observatory.mozilla.org/analyze/pixivfe.exozy.me) | ![Status](https://img.shields.io/website?url=https%3A%2F%2Fpixivfe.exozy.me&label=status)    |
| dragoongoose      | [https://pixivfe.drgns.space](https://pixivfe.drgns.space) | US      | No          | [![Mozilla HTTP Observatory Grade](https://img.shields.io/mozilla-observatory/grade-score/pixivfe.drgns.space?label=)](https://observatory.mozilla.org/analyze/pixivfe.drgns.space) | ![Status](https://img.shields.io/website?url=https%3A%2F%2Fpixivfe.drgns.space&label=status) |
| ducks.party       | [https://pixivfe.ducks.party](https://pixivfe.ducks.party) | NL      | No          | [![Mozilla HTTP Observatory Grade](https://img.shields.io/mozilla-observatory/grade-score/pixivfe.ducks.party?label=)](https://observatory.mozilla.org/analyze/pixivfe.ducks.party) | ![Status](https://img.shields.io/website?url=https%3A%2F%2Fpixivfe.ducks.party&label=status) |
| perennialte.ch    | [https://pixiv.perennialte.ch](https://pixiv.perennialte.ch) | AU      | No          | [![Mozilla HTTP Observatory Grade](https://img.shields.io/mozilla-observatory/grade-score/pixiv.perennialte.ch?label=)](https://observatory.mozilla.org/analyze/pixiv.perennialte.ch) | ![Status](https://img.shields.io/website?url=https%3A%2F%2Fpixiv.perennialte.ch&label=status)|
| darkness.services | [https://pixivfe.darkness.services](https://pixivfe.darkness.services) | US      | Yes         | [![Mozilla HTTP Observatory Grade](https://img.shields.io/mozilla-observatory/grade-score/pixivfe.darkness.services?label=)](https://observatory.mozilla.org/analyze/pixivfe.darkness.services) | ![Status](https://img.shields.io/website?url=https%3A%2F%2Fpixivfe.darkness.services&label=status) |
| thebunny.zone     | [https://pixivfe.thebunny.zone](https://pixivfe.darkness.services) | HR      | No         | [![Mozilla HTTP Observatory Grade](https://img.shields.io/mozilla-observatory/grade-score/pixivfe.thebunny.zone?label=)](https://observatory.mozilla.org/analyze/pixivfe.thebunny.zone) | ![Status](https://img.shields.io/website?url=https%3A%2F%2Fpixivfe.thebunny.zone&label=status) |
| thebunny.zone (🧅)| [http://pixivfe.bunny5exbgbp4sqe2h2rfq2brgrx3dhohdweonepzwfgumfyygb35wyd.onion](http://pixivfe.bunny5exbgbp4sqe2h2rfq2brgrx3dhohdweonepzwfgumfyygb35wyd.onion/) | HR      | No         | [![Mozilla HTTP Observatory Grade](https://img.shields.io/mozilla-observatory/grade-score/pixivfe.thebunny.zone?label=)](https://observatory.mozilla.org/analyze/pixivfe.thebunny.zone) | ![Status](https://img.shields.io/website?url=https%3A%2F%2Fpixivfe.thebunny.zone&label=status) |

If you are hosting your own instance, you can create a pull request to add it here!

For more information on instance uptime, see the [PixivFE instance status page](https://stats.uptimerobot.com/FbEGewWlbX).

## License

License: [AGPL3](https://www.gnu.org/licenses/agpl-3.0.txt)
