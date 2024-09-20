# Guidelines

This file contains guidelines for code development.

## Repository structure

- `assets/`: Static files (images, CSS, JS) and site templates.
- `config/`: Handle configurations for the server.
  - `config.go`: Parses for environment variables and pass down the configurations.
  - `proxy_list.go`: Contains the built-in third party image proxy list.
- `core/`: Make requests to Pixiv's server and parses information into structured data.
- `doc/`: Contains documentations for the general user and developers.
- `server/`: The web server.
  - `audit/`: Audit requests made by PixivFE to external APIs.
  - `handlers/`: Contains middlewares for logging, rate limiting, error handling, etc.
  - `proxy_checker/`: The image proxy checker.
  - `routes/`: Routes handling. Pages will be rendered here.
  - `session/`: Handles user's options.
  - `template/`: Core template renderer and template functions.
  - `utils/`: Other utilities.

## Naming
### Files

- Directories names **must** use [snake_case](https://en.wikipedia.org/wiki/Snake_case). 
- File names **must** use [snake-case](https://en.wikipedia.org/wiki/Snake_case), but with a dash (`-`) as a replacement for spaces (` `)
- All characters in directories and file names **must** be lowercase.

### Code

#### Variables
- Local variables should be named using [camelCase](https://en.wikipedia.org/wiki/Camel_case) with an initial lowercase letter.
- Global variables **must** be named using [CamelCase](https://en.wikipedia.org/wiki/Camel_case) with an initial uppercase letter.
- Environment variables **must** be named using [SNAKE_CASE](https://en.wikipedia.org/wiki/Snake_case) and all characters **must** be uppercase.
- Variable names should not contain any special characters.

## Adding features

- Edit the `Makefile` if you want to create scripts that helps the process of developing and/or using PixivFE.
- Don't create/add any shell files and configuration files of external programs in the root directory unless necessary.
- Add comments where necessary.