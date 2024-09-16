# Hosting PixivFE

PixivFE can be installed using various methods. This guide covers installation using [Docker](#docker) (recommended for production) and using a binary with a Caddy reverse proxy.

!!! note
    PixivFE requires a Pixiv account token to access the API. Refer to the [Obtaining the `PIXIVFE_TOKEN` cookie](obtaining-pixivfe-token.md) guide for detailed instructions.

## Docker

[Docker](https://www.docker.com/) lets you run containerized applications. Containers are loosely isolated environments that are lightweight and contain everything needed to run the application, so there's no need to rely on what's installed on the host.

Docker images for PixivFE are available at [Docker Hub](https://hub.docker.com/r/vnpower/pixivfe), with support for the amd64 platform.

### Docker Compose

!!! warning
    Deploying PixivFE using Docker Compose requires the Compose plugin to be installed. Follow these [instructions on the Docker Docs](https://docs.docker.com/compose/install) on how to install it.

#### 1. Setting up the repository

Clone the PixivFE repository and navigate to the directory:

```bash
git clone https://codeberg.org/VnPower/PixivFE.git && cd PixivFE
```

#### 2. Configure environment variables

Copy `.env.example` to `.env` and configure the variables as needed. Refer to the [Environment variables](environment-variables.md) page for more information.

!!! note
    Ensure you set `PIXIVFE_HOST=0.0.0.0` in the `.env` file.
    
    This allows PixivFE to bind to all network interfaces inside the container, which is necessary for Docker's network management to function correctly. The network access restrictions will be handled by Docker itself, not within PixivFE.

#### 3. Set token

Set the `PIXIVFE_TOKEN` environment variable in your `.env` file. This should be the value of the `PHPSESSID` cookie from your Pixiv account. For detailed instructions on obtaining this token, refer to the [Obtaining the `PIXIVFE_TOKEN` cookie](obtaining-pixivfe-token.md) guide.

#### 4. Compose!

Run `docker compose up -d` to start PixivFE.

To view the container logs, run `docker logs -f pixivfe`.

### Docker CLI

#### 1. Setting up the repository

Clone the PixivFE repository and navigate to the directory:

```bash
git clone https://codeberg.org/VnPower/PixivFE.git && cd PixivFE
```

#### 2. Configure environment variables

Copy `.env.example` to `.env` and configure the variables as needed. Refer to the [Environment variables](environment-variables.md) page for more information.

!!! note
    Ensure you set `PIXIVFE_HOST=0.0.0.0` in the `.env` file.
    
    This allows PixivFE to bind to all network interfaces inside the container, which is necessary for Docker's network management to function correctly. The network access restrictions will be handled by Docker itself, not within PixivFE.

#### 3. Deploying PixivFE

Run the following command to deploy PixivFE:

=== "Default port (`8282`)"

    ```bash
    docker run -d --name pixivfe -p 8282:8282 --env-file .env vnpower/pixivfe:latest
    ```

=== "Custom port (e.g., `8080`)"

    ```bash
    docker run -d --name pixivfe -p 8080:8282 --env-file .env vnpower/pixivfe:latest
    ```

If you're planning to use a reverse proxy, modify the port binding to only listen on the localhost port (e.g., `127.0.0.1:8282:8282`). This ensures that PixivFE listens only on the localhost, making it accessible solely through the reverse proxy.

## Binary

This setup uses [Caddy](https://caddyserver.com/) as the reverse proxy. Caddy is a great alternative to [NGINX](https://nginx.org/en/) because it is written in the [Go programming language](https://go.dev/), making it more lightweight and efficient. Additionally, Caddy is easy to configure, providing a simple and straightforward way to set up a reverse proxy.

### 1. Setting up the repository

Clone the PixivFE repository and navigate to the directory:

```bash
git clone https://codeberg.org/VnPower/PixivFE.git && cd PixivFE
```

### 2. Configure environment variables

Copy `.env.example` to `.env` and configure the variables as needed. Refer to the [Environment variables](environment-variables.md) page for more information.

### 3. Building and running PixivFE

PixivFE uses a [Makefile](https://www.gnu.org/software/make/manual/make.html#Introduction) to simplify the build and run process.

To build and run PixivFE, use the following commands:

```bash
make build
make run
```

This will build the PixivFE binary and start it. It will be accessible at `localhost:8282`.

### 4. Deploying Caddy

[Install Caddy](https://caddyserver.com/docs/install) using your package manager.

In the PixivFE directory, create a file named `Caddyfile` with the following content:

```caddy
example.com {
  reverse_proxy localhost:8282
}
```

Replace `example.com` with your domain and `8282` with the PixivFE port if you changed it.

Run `caddy run` to start Caddy.

## Updating

To update PixivFE to the latest version, follow the steps below that are relevant to your deployment method.

### Docker

#### Docker Compose

1. Pull the latest Docker image and repository changes:
   ```bash
   docker compose pull && git pull
   ```

2. Restart the container:
   ```bash
   docker compose up -d
   ```

#### Docker CLI

1. Pull the latest Docker image and repository changes:
   ```bash
   docker pull vnpower/pixivfe:latest && git pull
   ```

2. Stop and remove the existing container:
   ```bash
   docker stop pixivfe && docker rm pixivfe
   ```

3. Restart the container:
   ```bash
   docker run -d --name pixivfe -p 8282:8282 --env-file .env vnpower/pixivfe:latest
   ```

### Binary

1. Pull the latest changes from the repository:
   ```bash
   git pull
   ```

2. Rebuild and start PixivFE:
   ```bash
   make build
   make run
   ```

## Acknowledgements

- [Keep Caddy Running](https://caddyserver.com/docs/running#keep-caddy-running)
