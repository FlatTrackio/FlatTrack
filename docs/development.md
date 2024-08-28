# Development

Requirements:

-   [golang](https://golang.org/doc/install)
-   [nodejs](https://nodejs.org/en/download/)
-   a [Postgres](https://www.postgresql.org/download/) database

## Install dependencies with Brew

```sh
brew bundle install --file ./hack/Brewfile
```


## Write an environment file

Write the credientials for the database into `.env`

    APP_DB_USERNAME=flattrack
    APP_DB_PASSWORD=flattrack
    APP_DB_DATABASE=flattrack
    APP_DB_HOST=localhost


see the `example.env` file in the root of the repo for more.


## Launch a workspace with Zellij

launch the workspace

```sh
./hack/start-dev.sh
```


## Set up Postgres

Docker/Podman:

    docker run -d \
           --name flattrack-postgres \
           -p 5432:5432 \
           -e POSTGRES_DB=flattrack \
           -e POSTGRES_USER=flattrack \
           -e POSTGRES_PASSWORD=flattrack \
           postgres:16.1-alpine3.19

Manually connecting to the Postgres database:

    psql -U flattrack -d flattrack -w -h localhost

(entering the password as FlatTrack)


## Frontend build

Navigate to the `web` folder.

Install Frontend dependencies

    npm i

Build the frontend

    npm run build


## Backend build

Navigate to the root of the repo.

Build the backend

    ./hack/verify-build-smoketest.sh


### Start the backend

    go run .


### Remove migrations

    gomigrate -source "file://$(pwd)/migrations" -database postgres://flattrack:flattrack@localhost/flattrack?sslmode=disable down


## Project structure and details


### API

Written in golang, the API is located in [pkg/routes](https://gitlab.com/flattrack/flattrack/-/tree/master/pkg/routes).

The features and areas are separated into packages.

1.  Testing

    Tests are located in [test/backend/e2e](https://gitlab.com/flattrack/flattrack/-/tree/master/test/backend/e2e). So far there are only e2e tests for FlatTrack&rsquo;s API.


### Frontend

Written in Vue.js + JavaScript, the frontend is located in [web](https://gitlab.com/flattrack/flattrack/-/tree/master/test/frontend).

The frontend makes requests to the backend to perform actions.


### Database

The migrations and database structuring is located in [migrations](https://gitlab.com/flattrack/flattrack/-/tree/master/migrations).
Each table is created with [golang-migrate](https://github.com/golang-migrate/migrate).


### Assets

Images are located in [web/assets](https://gitlab.com/flattrack/flattrack/-/tree/master/web/assets), these are used throughout the project (such as in the frontend, and readme).


## Docs

To run the docs in development, use:

    mkdocs serve


## Making a release checklist

Things to do before making a release:

-   ensure docs represent the latest changes
-   ensure linting passes

