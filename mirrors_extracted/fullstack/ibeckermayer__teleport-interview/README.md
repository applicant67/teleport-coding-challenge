# teleport-interview

Repository for the [Gravitational](https://goteleport.com/) [fullstack challenge](https://github.com/gravitational/careers/blob/master/challenges/fullstack/dashboard.pdf)

## Requirements and Certificate Setup

#### mkcert

This project relies on self-signed certificates, which can cause your browser to complain at you and other minor annoyances. A simple tool for solving this is [mkcert](https://github.com/FiloSottile/mkcert), which creates a local CA that your browser will trust and uses it to sign subsequently generated certificates. Install `mkcert`, and then navigate to the `certs/` directory and generate a certificate and key for the app by running

```bash
# cd certs/
mkcert -key-file localhost.key -cert-file localhost.crt 0.0.0.0 localhost 127.0.0.1 ::1
```

#### Docker

This project is built and tested with [Docker](https://www.docker.com/get-started) and [Docker Compose](https://docs.docker.com/compose/) (`20.10`/`1.27`)

## Development

Build the development Docker images by running

```bash
# Uncomment --no-cache below if you've made any changes to the relevant Dockerfile's
docker-compose -f docker-compose-dev.yml build #--no-cache
```

and run them by running

```bash
docker-compose -f docker-compose-dev.yml up
```

This will start a hot-reload webpack development server serving the React app running on port 8080 and a hot-reload go server running on port 8000. Access the app in the browser by navigating to [https://0.0.0.0:8080/](https://0.0.0.0:8080/).

## Production

Build the production Docker image by running

```bash
docker-compose -f docker-compose-prod.yml build
```

and run it by running

```bash
docker-compose -f docker-compose-prod.yml up
```

This will start a go server exposed over port 8000 serving the React app and go api. Access the app in the browser by navigating to [https://localhost:8000/](https://localhost:8000/).
