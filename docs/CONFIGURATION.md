- [Configurable values](#sec-1)

FlatTrack is solely configured to run with environment variables; Here are the variables to configure FlatTrack with.

# Configurable values<a id="sec-1"></a>

| Field                    | Description                                                               | Default        |
|------------------------ |------------------------------------------------------------------------- |-------------- |
| `APP_PORT`               | The port for web traffic to bind to                                       | `:8080` \*     |
| `APP_PORT_METRICS`       | The port for metrics traffic to bind to                                   | `:2112` \*     |
| `APP_PORT_HEALTH`        | The port for web traffic to bind to                                       | `:8081` \*     |
| `APP_DIST_FOLDER`        | The location of the frontend web assets                                   | `./dist`       |
| `APP_DB_HOST`            | The Postgres host to connect to                                           | `localhost`    |
| `APP_DB_USERNAME`        | The Postgres role to connect as                                           | `flattrack`    |
| `APP_DB_PASSWORD`        | The Postgres role's password to authenticate with                         | `flattrack`    |
| `APP_DB_DATABASE`        | The Postgres database to write data to                                    | `flattrack`    |
| `APP_DB_MIGRATIONS_PATH` | The path to the database migrations to migrate the Postgres database with | `./migrations` |

\*all interfaces
