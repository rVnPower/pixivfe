<!-- The indentation on this page is delicate; avoid auto-formatting it with tools that might break it -->

# Roadmap

This roadmap outlines the upcoming features and improvements for PixivFE. It provides an overview of what users can expect in future releases and what developers are currently working on or considering for implementation.

!!! tip

    Want to recommend a feature? [Create an issue on the PixivFE repository](https://codeberg.org/VnPower/PixivFE/issues/new).

## Table

Upcoming features will be assigned with a release version.
A release version will be published when all features assigned with that version have been implemented.

**Developers note:** Sort features by their release number. Only core developers should be allowed to assign release versions to features.

**Changelog:**
- 05/10/24: Pixivision and full landing page has been postponed for the new template rewrite.

| Features                                                  | Release |
|-----------------------------------------------------------|---------|
| AI/R15/R18/R18-G filtering                                | 2.9     |
| Novel series                                              | 2.9     |
|     > We are here     (move me down on every release)     | current |
| Testing                                                   | 2.10    |
| Git version display                                       | 2.10    |
| Pixivision (articles, tags, categories, RSS)              | 2.11    |
| Full landing page (recommended users, trending tags, ...) | 2.11    |
| Manga series                                              |         |
| Complete novel content support (furigana, pages, ...)     |         |
| Pixiv Sketch                                              |         |
| Pixiv Idea (pixiv.net/idea)                               |         |
| Pixiv Request (pixiv.net/request)                         |         |
| User discovery                                            |         |
| Semi-popular artworks                                     |         |
| Localization (l10n)                                       |         |
| App API (mobile API) support                              |         |
| Native Ugoira support                                     |         |
| Search page / Search suggestions                          |         |
| Dynamic image gallery                                     |         |
| CSS reorganization / Theming                              |         |

## To implement

### `/settings/`

- [x] Merge login page with settings page
- [x] Persistence (http-only secure cookies)
- [ ] [User Settings](features/user-customization.md)

### `/novel/`

- [ ] [Novel support](features/novels.md)

    Might need some ideas for the reader's UI.
    Allow options for font size and family?
    Black and white backgrounds?
    Theme support?

### `/series/`

- [ ] **Manga series**

    Serialized web comics.

    Example: [Pixiv Manga Series](https://www.pixiv.net/user/13651304/series/171013)

- [ ] **Novel series**

### Independent features

- [x] **Multiple tokens support**

    Now you can do `PIXIVFE_TOKEN=TOKEN_A,TOKEN_B`

- [ ] **Pixivision**

    [Pixivision](https://www.pixivision.net/en)

    Pretty good to discover new artworks n stuff.

    Implement by parsing the webpage.

- [ ] **RSS support for Pixivision**

- [ ] **Search page**

    A page to do more extensive searching.

    Might require JavaScript for search recommendation, if wanted.

- [ ] **Full landing page**

    There are a lot of sections for the landing page. [Pixiv Landing](https://www.pixiv.net/ajax/top/illust)

    The artwork parsing part has already been implemented flawlessly.

    We only have to write the frontend code for those sections.

- [ ] **Various interesting pages from Pixiv.net:**

    - [Pixiv Idea](https://www.pixiv.net/idea/)
    - [Pixiv Request](https://www.pixiv.net/request)
    - [Pixiv Contest](https://www.pixiv.net/contest/) (no AJAX endpoints)

## To consider

- **Speculative Fetching**

    Fetch images from pixiv and cache them while we send the response page to users. When they ask for those images, we already have those.

    **Reference:** [Caching](features/caching.md)

- **App API support**

    May be painful to implement.
    Required to fully replace Pixiv, if user actions won't work universally.

    **Reference:** [#7](https://codeberg.org/VnPower/PixivFE/issues/7)

- **User discovery**

    For discovery page.
    Pretty useless if user actions (following) don't work.

- **"Popular" artworks**

    Check the README of this repository:

    [Mashiro GitHub Repository](https://github.com/kokseen1/Mashiro)

- **i18n**

    The last thing to work on, probably.

## Misc

- [x] **Ranking page**

    A lot of options weren't implemented.

- [x] **Revisit ranking calendar**

    There should be a way to display R18 thumbnails now?
