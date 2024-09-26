---
hide:
  - navigation
---

# Instance list

**Warning: Any public instance that isn't in this list is considered untrustworthy. Use them at your own risk.**

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

# Instance rules

(This instance rules was originally written by [the Invidious project](https://github.com/iv-org/documentation/blob/master/docs/instances.md))

For an instance to be added to this list, it must comply with all of the rules listed below:

1. Instances MUST have been up for at least a week before it can be added to this list.
2. Instances MUST not be more than a month out of date compared with either the latest commit or latest release. Any instance that is more than a month out of date is considered unmaintained and will be removed from the list.
3. Instances MUST have a stable uptime of at least 80% ([according to UptimeRobot](https://stats.uptimerobot.com/FbEGewWlbX)).
4. Instances MUST be served via domain name.
5. Instances MUST be served via HTTPS (or/and onion).
6. Instances using any DDoS Protection / MITM MUST be marked as such (e.g. Cloudflare, DDoS-Guard).
7. Instances using any type of anti-bot protection MUST be marked as such.
8. Instances MUST NOT use any type of analytics, including external scripts of any kind.
9. Any system whose goal is to modify the content served to the user (i.e web server HTML rewrite) is considered the same as modifying the source code.
10. Instances running a modified source code:
    - MUST respect the [GNU AGPL](https://en.wikipedia.org/wiki/GNU_Affero_General_Public_License) by publishing their source code and stating their changes **before** they are added to the list
    - MUST publish any later modification in a timely manner
    - MUST contain a link to both the modified and original source code of Invidious in the footer.
11. Instances MUST NOT serve ads (sponsorship links in the banner are considered ads) NOR promote products.
12. Instances MUST NOT restrict or disallow the access / usage to any [natural person](https://en.wikipedia.org/wiki/Natural_person) (e.g. a country's IP range MUST NOT be blocked, access by a natural person MUST NOT be disallowed for arbitrary reason) - this rule doesn't apply to [juridical persons](https://en.wikipedia.org/wiki/Juridical_person).

**NOTE:** We reserve the right to decline any instance from being added to the list, and to remove / ban any instance breaking the aforementioned rules.
