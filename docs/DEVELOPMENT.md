- [Development cycle](#sec-1)
  - [(almost) Automatic](#sec-1-1)
  - [Manual](#sec-1-2)
    - [Bring up the deployments](#sec-1-2-1)
    - [Port-forward the database connection](#sec-1-2-2)
    - [Backend](#sec-1-2-3)
    - [Frontend](#sec-1-2-4)
- [Additional](#sec-2)
    - [Manually connecting to the Postgres database](#sec-2-0-1)
    - [Remove migrations](#sec-2-0-2)
- [Project structure and details](#sec-3)
    - [Code](#sec-3-0-1)
    - [API](#sec-3-0-2)
    - [Frontend](#sec-3-0-3)
    - [Database](#sec-3-0-4)
    - [Assets](#sec-3-0-5)
- [Docs](#sec-4)
- [Making a release checklist](#sec-5)


# Development cycle<a id="sec-1"></a>

In-cluster local development is recommended, use [minikube](https://minikube.sigs.k8s.io) or [kind](https://kind.sigs.k8s.io/).

## (almost) Automatic<a id="sec-1-1"></a>

Get tilt from [Tilt.dev](https://tilt.dev).

```shell
tilt up --host 0.0.0.0
```

Build the backend:

```shell
CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags "-extldflags '-static' -s -w" -o flattrack ./main.go
```

Build the frontend:

```shell
npm run build
```

Notes:

-   this is the quickest development cycle factoring:
    -   multi-stage build
    -   components
    -   setup

## Manual<a id="sec-1-2"></a>

### Bring up the deployments<a id="sec-1-2-1"></a>

```shell
kubectl apply -k deployments/k8s-manifests/development/postgres
```

### Port-forward the database connection<a id="sec-1-2-2"></a>

```shell
kubectl -n flattrack-dev port-forward service/postgres 5432:5432
```

### Backend<a id="sec-1-2-3"></a>

```shell
go build -o flattrack ./main.go && ./flattrack
```

### Frontend<a id="sec-1-2-4"></a>

Install frontend dependencies:

```shell
npm i
```

Build the frontend:

```shell
npm run build
```

# Additional<a id="sec-2"></a>

### Manually connecting to the Postgres database<a id="sec-2-0-1"></a>

```shell
kubectl -n flattrack-dev exec -it postgres-0 -- psql
```

### Remove migrations<a id="sec-2-0-2"></a>

```shell
gomigrate -source "file://$(pwd)/migrations" -database postgres://flattrack:flattrack@localhost/flattrack?sslmode=disable down
```

# Project structure and details<a id="sec-3"></a>

### Code<a id="sec-3-0-1"></a>

### API<a id="sec-3-0-2"></a>

Written in golang, the API is located in [pkg/routes](https://gitlab.com/flattrack/flattrack/-/tree/master/pkg/routes).

The features and areas are separated into packages.

1.  Testing

    Tests are located in [test/backend/e2e](https://gitlab.com/flattrack/flattrack/-/tree/master/test/backend/e2e). So far there are only e2e tests for FlatTrack's API.

### Frontend<a id="sec-3-0-3"></a>

Written in Vue.js + JavaScript, the frontend is located in [web](https://gitlab.com/flattrack/flattrack/-/tree/master/test/frontend).

The frontend makes requests to the backend to perform actions.

### Database<a id="sec-3-0-4"></a>

The migrations and database structuring is located in [migrations](https://gitlab.com/flattrack/flattrack/-/tree/master/migrations). Each table is created with [golang-migrate](https://github.com/golang-migrate/migrate).

### Assets<a id="sec-3-0-5"></a>

Images are located in [web/assets](https://gitlab.com/flattrack/flattrack/-/tree/master/web/assets), these are used throughout the project (such as in the frontend, and readme).

# Docs<a id="sec-4"></a>

To run the docs in development, use:

```sh
docker run --rm -it -p 8000:8000 -v ${PWD}:/docs:ro,Z squidfunk/mkdocs-material
```

# Making a release checklist<a id="sec-5"></a>

Things to do before making a release:

-   update helm chart version
-   ensure docs represent the latest changes
-   ensure linting passes
