### mailx-google-service

`mailx-google-service` is an open-source microservice that consumes the [Gmail API](https://developers.google.com/gmail/api) for an open-souce email client *[Mailx](https://developers.google.com/gmail/api)*

### mailx-google-service tools
`mailx-google-service` uses multiple libraries and tools to work such as:
- `gorila mux` for a multiplexer router.
- `go-kit` standard library for microservices architecture.
- `sqlx` an extension for `database/sql`golang package.
- `goose` for incremental or decremental migrations.
- `pq` driver for a `postgres` database.
- `oauth2` for authorization.

### Getting started

To get started with `mailx-google-service` is important to have the following tools installed in your machine:
- `go v1.17+` (primary programming language)
- `Docker`
- `goose` (migration tool)

Clone the project in your desired location such as on `$GOPATH`:
```sh
git clone https://github.com/orlandorode97/mailx-google-service.git
```

`mailx-google-services` requires a `.env` at the root of the project with the following variables:
```sh
POSTGRES_HOST=
POSTGRES_PORT=
POSTGRES_USER=
POSTGRES_PASSWORD=
POSTGRES_DB_NAME=
POSTGRES_SSL_MODE=

GOOGLE_REDIRECT_URL=
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
```
After setting up the `.env` at the root the project run:
```sh
make build-run
```
The previous command builds a `mailx-google-service` container along with another container for a `postgres` database. Check the [Makefile](https://github.com/orlandorode97/mailx-google-service/blob/main/Makefile) for more available commands.

### Migrations
To perform migrations it's required to have installed `goose` in your machine. Check [tools section](#Getting-started).
To create a migration file run the following command:
```sh
goose create migration_file_name
```
The previous command will generate a go file where the needed sql statements will be placed. This migration has the following format `yyyymmddhhmmss_migration_name.go` and is created at the root of the project, you should move the migration file into the `migrations` folder.

When the container `mailx-google-service` is up and running run the follwing command to run a migration:
```sh
make goose-up
```

### TODO
A lot of things 😳
