# PixivFE

A privacy-respecting alternative front-end for Pixiv that doesn't suck.

<p>
<a href="https://codeberg.org/vnpower/pixivfe">
<img alt="Get it on Codeberg" src="https://get-it-on.codeberg.org/get-it-on-blue-on-white.png" height="60">
</a>
</p>

[![CI badge](https://ci.codeberg.org/api/badges/12556/status.svg)](https://ci.codeberg.org/repos/12556)
[![Go Report Card](https://goreportcard.com/badge/codeberg.org/vnpower/pixivfe/v2)](https://goreportcard.com/report/codeberg.org/vnpower/pixivfe)

Questions? Feedback? You can [PM me](https://matrix.to/#/@vnpower:eientei.org) on Matrix! You can also see the [Known quirks](https://pixivfe-docs.pages.dev/known-quirks/) page to check if your issue has a known solution.

You can keep track of this project's development using the [roadmap](doc/dev/general.md).

## Features

- Lightweight - both the interface and the code
- Privacy-first - the server will do the work for you
- No bloat - we only serve HTML, CSS and minimal JS code
- Open source - you can trust me!

## Development

Prerequisites:

- [Go](https://go.dev/doc/install) (to build PixivFE from source)
- [Sass](https://github.com/sass/dart-sass/) (will be run by PixivFE in development mode)
- [Make](https://www.gnu.org/software/make/) (to use the Makefile)

To install Dart Sass, you can choose any of the following methods.

- use system package manager (usually called `dart-sass`)
- download executable from [the official release page](https://github.com/sass/dart-sass/releases)
- `pnpm i -g sass`

Then, run the project:

```bash
# Clone the PixivFE repository
git clone https://codeberg.org/VnPower/PixivFE.git && cd PixivFE

# Build the project
make build

# Run PixivFE in development mode (styles and templates reload automatically)
PIXIVFE_DEV=1 <other_environment_variables> make run
```

## Using the Makefile

The project includes a Makefile to simplify common development tasks. Run `make help` to view the available commands.

## Hosting PixivFE

You can use PixivFE for personal use! Assuming that you use an operating system that can run POSIX shell scripts, install `go`, clone this repository, and use the Makefile to build and run the project.
I recommend self-hosting your own instance for personal use, instead of relying entirely on official instances.

To deploy PixivFE using Docker or the compiled binary, see [Hosting PixivFE](https://pixivfe-docs.pages.dev/hosting-pixivfe/).

### Public Instances

<!-- The current instance table is really wide; maybe there's a better way of formatting it without losing information?
The badges are also difficult to read on a small screen due to Codeberg shrinking the width of the columns -->

| Name              | URL                                             | Country | Cloudflare? | Observatory Grade                                                                                                                               | Status                                                                                                                                               |
|-------------------|-------------------------------------------------|---------|-------------|-------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------|
| exozyme (Official)| [https://pixivfe.exozy.me](https://pixivfe.exozy.me) | US      | No          | [![MDN HTTP Observatory Grade](https://img.shields.io/mozilla-observatory/grade-score/pixivfe.exozy.me?label=)](https://developer.mozilla.org/en-US/observatory/analyze?host=pixivfe.exozy.me) | ![Status](https://img.shields.io/uptimerobot/status/m796383741-c72f1ae6562dc943d032ba96)    |
| dragoongoose      | [https://pixivfe.drgns.space](https://pixivfe.drgns.space) | US      | No          | [![MDN HTTP Observatory Grade](https://img.shields.io/mozilla-observatory/grade-score/pixivfe.drgns.space?label=)](https://developer.mozilla.org/en-US/observatory/analyze?host=pixivfe.drgns.space) | ![Status](https://img.shields.io/uptimerobot/status/m796383743-c0cf0d6b5dbb09c8dbe7dc53) |
| ducks.party       | [https://pixivfe.ducks.party](https://pixivfe.ducks.party) | NL      | No          | [![MDN HTTP Observatory Grade](https://img.shields.io/mozilla-observatory/grade-score/pixivfe.ducks.party?label=)](https://developer.mozilla.org/en-US/observatory/analyze?host=pixivfe.ducks.party) | ![Status](https://img.shields.io/uptimerobot/status/m796383747-c92c281f520d52fe3fd894ed) |
| perennialte.ch    | [https://pixiv.perennialte.ch](https://pixiv.perennialte.ch) | AU      | No          | [![MDN HTTP Observatory Grade](https://img.shields.io/mozilla-observatory/grade-score/pixiv.perennialte.ch?label=)](https://developer.mozilla.org/en-US/observatory/analyze?host=pixiv.perennialte.ch) | ![Status](https://img.shields.io/uptimerobot/status/m796383748-503799f65873a23dbc860a02) |
| darkness.services | [https://pixivfe.darkness.services](https://pixivfe.darkness.services) | US      | Yes         | [![MDN HTTP Observatory Grade](https://img.shields.io/mozilla-observatory/grade-score/pixivfe.darkness.services?label=)](https://developer.mozilla.org/en-US/observatory/analyze?host=pixivfe.darkness.services) | ![Status](https://img.shields.io/uptimerobot/status/m796758268-211b0a18f07b88673820715f) |
| thebunny.zone     | [https://pixivfe.thebunny.zone](https://pixivfe.thebunny.zone) | HR      | No          | [![MDN HTTP Observatory Grade](https://img.shields.io/mozilla-observatory/grade-score/pixivfe.thebunny.zone?label=)](https://developer.mozilla.org/en-US/observatory/analyze?host=pixivfe.thebunny.zone) | ![Status](https://img.shields.io/uptimerobot/status/m797561997-78a2d28dadf458745d556322) |
| thebunny.zone (🧅)| [http://pixivfe.bunny5exbgbp4sqe2h2rfq2brgrx3dhohdweonepzwfgumfyygb35wyd.onion](http://pixivfe.bunny5exbgbp4sqe2h2rfq2brgrx3dhohdweonepzwfgumfyygb35wyd.onion/) | HR      | No          | [![MDN HTTP Observatory Grade](https://img.shields.io/mozilla-observatory/grade-score/pixivfe.thebunny.zone?label=)](https://developer.mozilla.org/en-US/observatory/analyze?host=pixivfe.thebunny.zone) | ![Status](https://img.shields.io/uptimerobot/status/m797561997-78a2d28dadf458745d556322) |

If you are hosting your own instance, you can create a pull request to add it here!

For more information on instance uptime, see the [PixivFE instance status page](https://stats.uptimerobot.com/FbEGewWlbX).

This information is duplicated at https://pixivfe-docs.pages.dev/instance-list/.

## Hosting Image Proxy Server

PixivFE can work with or without an external image proxy server. Here is the [list of public image proxies](https://pixivfe-docs.pages.dev/public-image-proxies/).
See [hosting a Pixiv image proxy](https://pixivfe-docs.pages.dev/hosting-image-proxy-server/) if you want to host one yourself.
