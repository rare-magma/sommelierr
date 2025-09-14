# sommelierr

Service that selects a random movie and/or series from Radarr / Sonarr.

## Dependencies

- [go](https://go.dev/)
- Optional:
  - [docker](https://docs.docker.com/)

## Relevant documentation

- [Radarr API](https://radarr.video/docs/api/)
- [Sonarr API](https://sonarr.tv/docs/api)

## Installation

### With Docker

#### docker-compose

1. Configure `.env` (see the configuration section below).
1. Run it.

   ```bash
   docker compose up --detach
   ```

#### docker build & run

1. Build the docker image.

   ```bash
   docker build . --tag sommelierr
   ```

1. Configure `.env` (see the configuration section below).
1. Run it.

   ```bash
    docker run --rm --read-only --env-file .env --publish 8080:8080 --cap-drop ALL --security-opt no-new-privileges:true --cpus 2 -m 64m --pids-limit 16 ghcr.io/rare-magma/sommelierr:latest
    ```

1. Open Sommelierr's UI @ <http://localhost:8080>

### Manually

1. Build it:

   ```bash
   go build -ldflags "-s -w" -trimpath -o sommelierr ./cmd/server
   ```

2. Run it:

   ```bash
   ./sommelierr
   ```

3. Open Sommelierr's UI @ <http://localhost:8080>

### Environment file

The env file has a few options:

```plaintext
EXCLUDE_LABEL=watched
RADARR_HOST=https://radarr.example.com
RADARR_API_KEY=aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
SONARR_HOST=https://sonarr.example.com
SONARR_API_KEY=bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
PORT=8080
```

- `EXCLUDE_LABEL` should be the label used to exclude already watched movies/series. This is optional.
- `RADARR_HOST` should be the FQDN of the Radarr service.
- `RADARR_API_KEY` should be the Radarr api key. This is configured in Settings - General - API Key
- `SONARR_HOST` should be the FQDN of the Sonarr service.
- `SONARR_API_KEY` should be the Sonarr api key. This is configured in Settings - General - API Key
- `PORT` should be the port available for Sommelierr to attach to. Defaults to 8080

## Troubleshooting

Check the service logs for errors with Radarr/Sonarr's API.

## Credits

- [Radarr](https://radarr.video/)
- [Sonarr](https://sonarr.tv/)
