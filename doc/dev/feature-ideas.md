# Proposals

This section contains some potential features or redesigns that can be implemented into PixivFE.

## Pixivision

**Summary**: Pixivision is a service owned by Pixiv that publishes articles about various types of artwork themes...

### Notes

- [Pixivision](https://www.pixivision.net/en/) is an independent service. Most of the data comes from Pixivision's website.
- Pixivision does not provide an API for data access. Pages on Pixivision seems to be static, so the HTML could be easily parsed.

### Thoughts

- Write a separate module for Pixivision and integrate it into PixivFE.
- We do web scraping and HTML parsing for this one.

## Sketch

**Summary**: Sketch is a service owned by Pixiv that allow users to livestream, mostly for their drawing process.

### Notes

- [Sketch](https://sketch.pixiv.net/) is an independent service. Most of the data comes from Pixivision's website.
- Sketch has a dedicated API for data access. Sketch uses the same type of authentication that Pixiv has.
- Detailed notes TBA.

### Thoughts

- Thanks to the public API, pages could be build easily.
- For the streaming part, we may have to include a JavaScript library for HLS streaming.

## Ugoira support

**Summary:** Ugoiras are Pixiv's "animated image" format.

### Notes

- Ugoiras are basically a bunch of sorted images combined with a fixed delay for each of them.
- Pixiv provides one JSON endpoint for delays and filenames and one endpoint for the (ZIP) images archive.
- One has to write their own player based on things Pixiv provide.
- You can check out Pixiv's implementation on their own ugoira player [here](https://github.com/pixiv/zip_player).

### Thoughts

- GIF/APNG/WEBP renderer.
- Some people want to convert ugoiras to video formats? (no idea)

## Landing page

**Summary**: PixivFE's homepage.

### Notes

- Pixiv's homepage contains a lot of interesting content.
- PixivFE's backend for the landing page already implemented almost all of the data from Pixiv.
- The only thing left is to write the frontend for them.
- Detailed notes TBA.

### Thoughts

- Spend some time to write some HTML/SCSS.
- Currently, you have to authenticate (login) in order to access the _full_ landing page. Can we show the _full_ page to unauthenticated users as well?

## Popular artworks

**Summary**: Pixiv has a "Sort by views" and "Sort by bookmarks" feature that is only available for premium users.

### Notes

- There are some search ["hacks"](https://github.com/kokseen1/Mashiro/) that could yield relatively accurate results for popular artworks.

### Ideas

- Look into repos that attempts to retrieve popular artworks
- If search "hacking" is possible, could there be more "hacks" around?

## "User discovery" page

**Summary**: Like artwork discovery, but it is for users.

### Notes

- Currently, we do not know if we could implement the "user follow" function into PixivFE.
- The development for this page has been put on hold because of it, since "following", after all, is what you want to do if you discover an user you like.

### Ideas

- It is easy to implement thanks to the API.

## Search suggestions

**Summary**: Pixiv provides [an API endpoint](https://www.pixiv.net/ajax/search/suggestion?mode=all&lang=en) for search suggestions.

### Notes

- The search suggestions appear when you focus on the search bar.

### Ideas

- We can prefetch the search suggestions for every request on PixivFE. But this means we will have to add one request (to Pixiv's API) for each PixivFE page request. ([Caching?](features/caching.md))
- We can implement JavaScript to fetch the suggestions every time the user focuses on the search bar.
- We can create a separate page just for this.

## App API support

**Summary**: Apart of the public AJAX API, Pixiv also provides a private API, used specifically for mobile applications.

### Notes

- Because you already could do almost everything through the AJAX API, there is really no point to integrate the App API.
- I added this section because there are some limitations to the public API (following,...).

### Ideas

- Write more stuff when desperate.

## Novel page

## Image grid layout

## Series

## Server's PixivFE Git version/commit

**Summary**: Implementation is essentially complete. However, the removal of `.dockerignore` as a dirty workaround[^1] is unfortunate.

### Ideas

- Implemented by other open-source alternative frontends, for example:
    - [Invidious](https://github.com/iv-org/invidious/blob/a021b93063f3956fc9bb3cce0fb56ea252422738/src/invidious/views/template.ecr#L117-L131)
    - [Nitter](https://github.com/zedeus/nitter/blob/b62d73dbd373f08af07c7a79efcd790d3bc1a49c/src/views/about.nim#L5-L9)

### Notes

- Initial implementation by [jackyzy823](https://codeberg.org/jackyzy823) in [#104](https://codeberg.org/VnPower/PixivFE/pulls/104) ([f53e1c3e4d](https://codeberg.org/VnPower/PixivFE/pulls/104/commits/f53e1c3e4db31587ede84f5518d729fcc076dd44)) via a `REVISION` variable defined as the Git commit hash of HEAD at build time.

- The `REVISION` variable was modified to include both the commit date and hash in [`901286d98e`](https://codeberg.org/VnPower/PixivFE/commit/901286d98ec27faa7f255146ce38d7c4a87f30ed).

- A "dirty" flag is appended to the `REVISION` variable if there are uncommitted changes in [`7a9216a165`](https://codeberg.org/VnPower/PixivFE/commit/7a9216a165a10fda24666e256747420f56473f0f).

- The `.dockerignore` file was removed to prevent Docker image builds from always being flagged as "dirty" in [`436f4073ea`](https://codeberg.org/VnPower/PixivFE/commit/436f4073eaf6168946674126fe61626ba3753afd).

## Download all buttons for all containers

## Page profile

## Download button in artwork page

## Option to select default image quality in artwork page

[^1]: No pun intended.
