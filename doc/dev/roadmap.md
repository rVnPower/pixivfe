# Roadmap

Planned features.

## Table

Upcoming features will be assigned with a release version. 
A release version will be published when all features assigned with that version have been implemented.

Want to recommend a feature? Create an issue [here](https://codeberg.org/VnPower/PixivFE/issues/new?).

**Developers note:** Sort features by their release number. Only core developers should be allowed to assign release versions to features.

| Features                                                  | Release |
|-----------------------------------------------------------|---------|
| Manga series                                              | N/A     |
| Novel series                                              | N/A     |
| Complete novel content support (furigana, pages, ...)     | N/A     |
| Pixivision (articles, tags, categories, RSS)              | N/A     |
| Full landing page (recommended users, trending tags, ...) | N/A     |
| Pixiv Sketch                                              | N/A     |
| Pixiv Idea (pixiv.net/idea)                               | N/A     |
| Pixiv Request (pixiv.net/request)                         | N/A     |
| User discovery                                            | N/A     |
| Semi-popular artworks                                     | N/A     |
| Localization (l10n)                                       | N/A     |
| App API (mobile API) support                              | N/A     |
| Native Ugoira support                                     | N/A     |
| Search page / Search suggestions                          | N/A     |
| Dynamic image gallery                                     | N/A     |
| CSS reorganization / Theming                              | N/A     |
| Git version display                                       | N/A     |
| AI/R15/R18/R18-G filtering                                | N/A     |

## To implement

/settings/

- [x] Merge login page with settings page
- [x] Persistence (http-only secure cookies)
- [User Settings](features/user-customization.md)

/novel/

- [Novel support](features/novels.md)  
Might need some ideas for the reader's UI.  
Allow options for font size and family?  
Black and white backgrounds?  
Theme support?  

/series/
- [ ] Manga series  
Serialized web comics. Example: https://www.pixiv.net/user/13651304/series/171013
- [ ] Novel series  

Independent features

- [x] Multiple tokens support
Now you can do PIXIVFE_TOKEN=TOKEN_A,TOKEN_B

- [ ] Pixivision  
https://www.pixivision.net/en/  
Pretty good to discover new artworks n stuff.  
Implement by parsing the webpage.

  - [ ] RSS support for Pixivision  

- [ ] Search page  
A page to do more extensive searching.  
Might require JavaScript for search recommendation, if wanted.


- [ ] Full landing page  
There are a lot of sections for the landing page. https://www.pixiv.net/ajax/top/illust  
The artwork parsing part has already been implemented flawlessly.  
We only have to write the frontend code for those sections.

- [ ] Various interesting pages from Pixiv.net  
  - https://www.pixiv.net/idea/
  - https://www.pixiv.net/request
  - https://www.pixiv.net/contest/ (no AJAX endpoints)

## To consider

- Speculative Fetching
Fetch images from pixiv and cache them while we send the response page to users. When they ask for those images, we already have those.

- App API support  
May be painful to implement.
Required to fully replace Pixiv, if user actions won't work universally.
https://codeberg.org/VnPower/PixivFE/issues/7

- User discovery  
For discovery page.  
Pretty useless if user actions (following) doesn't work.

- "Popular" artworks  
Check the README of this:  
https://github.com/kokseen1/Mashiro

- i18n  
The last thing to work on, probably.

## Misc

- [x] Ranking page  
A lot of options weren't implemented.

- [x] Revisit ranking calendar  
There should be a way to display R18 thumbnails now?
