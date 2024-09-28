# Pixiv App API

Besides the web API, Pixiv also has a private API made specifically for mobile applications.

There are currently no official API specifications for the private API, and the unofficial ones seem to lack several functionalities. We might have to write one.

Examples of implementations:
- https://github.com/upbit/pixivpy/blob/master/pixivpy3/aapi.py
- https://github.com/book000/pixivts/blob/master/src/pixiv.ts

## Goal
- Separate the Web API core and the App API core, but made both of them compatible with each other.
- Complications...

