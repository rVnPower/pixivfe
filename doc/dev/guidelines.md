This file contains guidelines for code development.

# File structure

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
