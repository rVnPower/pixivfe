# Instance list

This page lists all currently running instances of PixivFE that are available for use. They are ordered from oldest to newest based on when they were added. **Please scroll horizontally to see all columns.**

To check the uptime history and status of these instances, visit the [PixivFE instance status page](https://stats.uptimerobot.com/FbEGewWlbX).

Instances marked as having analytics load external tracking scripts, such as [Cloudflare Web Analytics](https://developers.cloudflare.com/analytics/web-analytics/).

!!! tip
    To add your instance to this list, [create an issue on the PixivFE repository](https://codeberg.org/VnPower/PixivFE/issues/new?template=.forgejo%2fissue_template%2fadd-instance.yaml) using the "Add Instance" template.

<!-- Note to page editors: The tables below only refresh their data when `mkdocs serve` is restarted, due to how the data is templated in from the CSV files.  -->

## Clearnet

These instances can be accessed through any regular web browser without any special setup.

{{ read_csv('data/instances.csv') }}

<!-- Human-readable list when viewing raw:

- Name: exozyme (Official)
  URL: https://pixivfe.exozy.me
  Country: US
  Cloudflare proxy: No
  Analytics: No

- Name: dragongoose
  URL: https://pixivfe.drgns.space
  Country: US
  Cloudflare proxy: No
  Analytics: No

- Name: ducks.party
  URL: https://pixivfe.ducks.party
  Country: NL
  Cloudflare proxy: No
  Analytics: No

- Name: perennialte.ch
  URL: https://pixiv.perennialte.ch
  Country: AU
  Cloudflare proxy: No
  Analytics: No

- Name: darkness.services
  URL: https://pixivfe.darkness.services
  Country: US
  Cloudflare proxy: Yes
  Analytics: No

- Name: thebunny.zone
  URL: https://pixivfe.thebunny.zone
  Country: HR
  Cloudflare proxy: No
  Analytics: No -->

## Tor onion services

These instances are only accessible using the Tor browser.

Since they are hosted on the Tor network, these instances provide better privacy compared to clearnet instances. However, they may have slower performance due to how onion routing works.

{{ read_csv('data/instances-onion.csv') }}

<!-- Human-readable list when viewing raw:

- Name: thebunny.zone
  URL: http://pixivfe.bunny5exbgbp4sqe2h2rfq2brgrx3dhohdweonepzwfgumfyygb35wyd.onion -->
