# A unified documentation for PixivFE

## Proposals

This section contains some potential features or redesigns that can be implemented into PixivFE.

### Pixivision
**Summary**: Pixivision is a service owned by Pixiv that publishes articles about various types of artwork themes...

**Notes**: [Pixivision](https://www.pixivision.net/en/) is an independent service. Most of the data comes from Pixivision's website.
Pixivision does not provide an API for data access. Pages on Pixivision seems to be static, so the HTML could be easily parsed.

**Thoughts**:
- Write a separate module for Pixivision and integrate it into PixivFE.
- We do web scraping and HTML parsing for this one.

### Sketch

**Summary**: Sketch is a service owned by Pixiv that allow users to livestream, mostly for their drawing process.

**Notes**: [Sketch](https://sketch.pixiv.net/) is an independent service. Most of the data comes from Pixivision's website.
Sketch has a dedicated API for data access. Sketch uses the same type of authentication that Pixiv has.
Detailed notes TBA.

**Thoughts**:
- Thanks to the public API, pages could be build easily.
- For the streaming part, we may have to include a JavaScript library for HLS streaming.

### Ugoira support

**Summary**: Ugoiras are Pixiv's "animated image" format.

**Notes**: Ugoiras are basically a bunch of sorted images combined with a fixed delay for each of them. Pixiv provides one JSON endpoint for delays and filenames
and one endpoint for the (ZIP) images archive.
One has to write their own player based on things Pixiv provide.
You can checkout Pixiv's implementation on their own ugoira player [here](https://github.com/pixiv/zip_player).

**Thoughts**:
- GIF/APNG/WEBP renderer.
- Some people want to convert ugoiras to video formats? (no idea)

### Landing page

**Summary**: PixivFE's homepage.

**Notes**: Pixiv's homepage contains a lot of interesting contents.
PixivFE's backend for the landing page already implemented almost all of the data from Pixiv.
The only thing left is to write the frontend for them. Detailed notes TBA.

**Thoughts**:
- Spend some time to write some HTML/SCSS.
- Currently, you have to authenticate (login) in order to access the *full* landing page. Can we show the *full* page to unauthenticated users as well?

### Popular artworks
### "User discovery" page
### Search page
### Novel page
### Image grid layout
### App API support
### Series
### Server's PixivFE Git version/commit
### Artwork filters

## Flaws

This section documents some bad/buggy designs in PixivFE's design, both frontend and backend.

### Cookies management / validation
### "Switchers"
### Undocumented code
### "kmutex"
### Browser compatibility / Universal CSS

## References

This section contains multiple external links to materials/resources that could help.

- [pixiv.pics](https://www.pixiv.pics/)
