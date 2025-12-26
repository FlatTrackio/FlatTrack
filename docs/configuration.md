# Configuration

FlatTrack is solely configured to run with environment variables (or a `.env` file in the same directory as the FlatTrack binary, such as `/app/.env`); Here are the variables to configure FlatTrack with.

## Configurable values

| Field                           | Description                                                                                                                   | Default               |
|---------------------------------|-------------------------------------------------------------------------------------------------------------------------------|-----------------------|
| `APP_ENV_FILE`                  | A `.env` file to load                                                                                                         | `.env`                |
| `APP_PORT` or `PORT`            | The port for web traffic to bind to                                                                                           | `:8080` *             |
| `APP_PORT_METRICS`              | The port for metrics traffic to bind to                                                                                       | `:2112` *             |
| `APP_PORT_HEALTH`               | The port for web traffic to bind to                                                                                           | `:8081` *             |
| `APP_WEB_FOLDER`                | The location of the frontend web assets                                                                                       | `./kodata/web`        |
| `APP_URL`                       | The location of the site                                                                                                      | `""`                  |
| `APP_DB_HOST`                   | The Postgres host to connect to                                                                                               | `localhost`           |
| `APP_DB_PORT`                   | The Postgres port use                                                                                                         | `5432`                |
| `APP_DB_USERNAME`               | The Postgres role to connect as                                                                                               | `flattrack`           |
| `APP_DB_PASSWORD`               | The Postgres role's password to authenticate with                                                                             | `flattrack`           |
| `APP_DB_DATABASE`               | The Postgres database to write data to                                                                                        | `flattrack`           |
| `APP_DB_SSLMODE`                | The Postgres SSL mode to use                                                                                                  | `disable`             |
| `APP_DB_CONNECTION_STRING`      | a full Postgres database connection string                                                                                    |                       |
| `APP_DB_MIGRATIONS_PATH`        | The path to the database migrations to migrate the Postgres database with                                                     | `./kodata/migrations` |
| `APP_METRICS_ENABLED`           | Serve Prometheus metrics endpoint                                                                                             | `"true"`              |
| `APP_MAINTENANCE_MODE`          | Set instance into a maintenance mode, disallowing access                                                                      |                       |
| `APP_HEALTH_ENABLED`            | Serve healthz endpoint                                                                                                        | `"true"`              |
| `APP_SMTP_ENABLED`              | Send SMTP email notifcations                                                                                                  | `"false"`             |
| `APP_SMTP_USERNAME`             | Username for SMTP emails                                                                                                      | `""`                  |
| `APP_SMTP_PASSWORD`             | Password for SMTP emails                                                                                                      | `""`                  |
| `APP_SMTP_HOST`                 | Host for SMTP emails                                                                                                          | `""`                  |
| `APP_SMTP_PORT`                 | Port for SMTP emails                                                                                                          | `""`                  |
| `APP_MINIO_ACCESS_KEY`          | The access key for a Minio storage bucket                                                                                     |                       |
| `APP_MINIO_SECRET_KEY`          | The secret key for a Minio storage bucket                                                                                     |                       |
| `APP_MINIO_BUCKET`              | The Minio storage bucket to use                                                                                               |                       |
| `APP_MINIO_HOST`                | The Minio host to use for storage                                                                                             |                       |
| `APP_MINIO_USE_SSL`             | Require SSL for a Minio connection                                                                                            |                       |
| `APP_HTTP_REAL_IP_HEADER`       | The header to use for forwarding IP address real value                                                                        | `""`                  |
| `APP_SETUP_MESSAGE`             | Display a message on the setup page                                                                                           | `""`                  |
| `APP_LOGIN_MESSAGE`             | Display a message on the login page                                                                                           | `""`                  |
| `APP_EMBEDDED_HTML`             | Add custom HTML to the head of index.html                                                                                     | `""`                  |
| `APP_REGISTRATION_SECRET`       | Require a matching registration secret to be passed during registration as the parameter `secret`                             |                       |
| `APP_SCHEDULER_USE_ENDPOINT`    | Use endpoint with scheduler at `/api/system/scheduler`                                                                        | `false`               |
| `APP_SCHEDULER_ENDPOINT_SECRET` | Set a secret for scheduler endpoint which must match header `X-FlatTrack-Scheduler-Secret` (required when scheduler disabled) |                       |
| `APP_LOG_LEVEL`                 | Sets the log level, between `INFO`, `DEBUG`, `WARN` and `ERROR`                                                               | `INFO`                |
| `APP_LOG_TIMEZONE`              | Sets the timezone for the logs. Defaults to UTC                                                                               |                       |
| `APP_MAINTENANCE_MODE_MESSAGE`  | A custom message to display when in maintenance mode                                                                          | `""`                  |


\*this port runs on all available interfaces - this may be better configured on a single host as `127.0.0.1:...` for security

