Dockerizing skillKonnect

What I added

- `Dockerfile` - multi-stage build that compiles the Go app with CGO enabled (needed for sqlite) and produces a small runtime image (Debian slim with libsqlite3)
- `docker-compose.yml` - runs the API service, exposes port 8080 on the host and mounts a named volume `app_data` for the sqlite DB file
- `.dockerignore` - exclude unnecessary files from the build context

How it works

- The container expects to use sqlite (DB_DRIVER=sqlite3) and DB_NAME is set to `/data/app_db` inside the container. The compose file mounts an `app_data` volume at `/data` so the sqlite DB file is persisted outside the container.
- The container's `HTTP_LISTEN_ADDR` is set to `:8080`, and the image exposes port 8080. Docker-compose maps host 8080 to container 8080.

Build & run locally (dev machine)

1. Build the image locally:

   docker compose build

2. Start the service:

   docker compose up -d

3. Confirm it's listening on host port 8080:

   curl -v http://localhost:8080/health || curl -v http://localhost:8080/

Notes & recommended follow-ups

- Migrations: the image does not automatically run DB migrations. You can run migrations inside the container (or extend the Dockerfile/entrypoint). For example:

  docker compose exec api /app/skillkonnect migrate # if you add migration logic

  Or use the goose command used in the Makefile with GO commands.

- Secrets: set `SUPERKIT_SECRET` in your environment or in a `.env` file used by docker-compose (don't commit secrets).
- If you prefer the app to listen directly on 8080 (instead of being proxied by templ), set `HTTP_LISTEN_ADDR=:8080` in your environment for local dev.

If you want, I can:

- Add a migration step to the container entrypoint and wire goose (or GORM auto-migrate) to run at startup
- Add a small health endpoint if the app doesn't already expose one
