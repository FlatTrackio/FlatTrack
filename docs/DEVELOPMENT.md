- [Development cycle](#sec-1)
  - [Bring up the deployments](#sec-1-1)
  - [Port-forward the database connection](#sec-1-2)
  - [Quick development patterns](#sec-1-3)
    - [Backend](#sec-1-3-1)
    - [Frontend](#sec-1-3-2)
  - [Auto build and launch with Tilt](#sec-1-4)
- [Additional](#sec-2)
  - [Manually connecting to the Postgres database](#sec-2-1)
  - [Remove migrations](#sec-2-2)
- [Project structure and details](#sec-3)
  - [Code](#sec-3-1)
    - [API](#sec-3-1-1)
    - [Frontend](#sec-3-1-2)
  - [Database](#sec-3-2)
  - [Assets](#sec-3-3)
- [Docs](#sec-4)


# Development cycle<a id="sec-1"></a>

In-cluster local development is recommended, use [minikube](https://minikube.sigs.k8s.io) or [kind](https://kind.sigs.k8s.io/).

## Bring up the deployments<a id="sec-1-1"></a>

```shell
kubectl apply -k k8s-manifests/development/postgres
```

## Port-forward the database connection<a id="sec-1-2"></a>

```shell
kubectl -n flattrack-dev port-forward service/postgres 5432:5432
```

## Quick development patterns<a id="sec-1-3"></a>

### Backend<a id="sec-1-3-1"></a>

```shell
go build -o flattrack src/backend/main.go && ./flattrack
```

### Frontend<a id="sec-1-3-2"></a>

```shell
npm i
```

```shell
npm run build:frontend
```

## Auto build and launch with Tilt<a id="sec-1-4"></a>

Get tilt from [Tilt.dev](https://tilt.dev).

```shell
tilt up --host 0.0.0.0
```

# Additional<a id="sec-2"></a>

## Manually connecting to the Postgres database<a id="sec-2-1"></a>

```shell
kubectl -n flattrack-dev exec -it postgres-0 -- psql
```

## Remove migrations<a id="sec-2-2"></a>

```shell
gomigrate -source "file://$(pwd)/migrations" -database postgres://postgres:postgres@localhost/flattrack?sslmode=disable down
```

# Project structure and details<a id="sec-3"></a>

## Code<a id="sec-3-1"></a>

### API<a id="sec-3-1-1"></a>

Written in golang, the API is located in [src/backend](../src/backend).

The features and areas are separated into packages.

1.  Testing

    Tests are located in [test/backend/e2e](../test/backend/e2e). So far there are only e2e tests for FlatTrack's API.

### Frontend<a id="sec-3-1-2"></a>

Written in Vue.js + JavaScript, the frontend is located in [src/frontend](../src/frontend).

The frontend makes requests to the backend to perform actions.

## Database<a id="sec-3-2"></a>

The migrations and database structuring is located in [migrations](../migrations). Each table is created with [golang-migrate](https://github.com/golang-migrate/migrate).

## Assets<a id="sec-3-3"></a>

Images are located in [src/frontend/assets](../src/frontend/assets), these are used throughout the project (such as in the frontend, and readme).

# Docs<a id="sec-4"></a>

To run the docs in development, use:

```sh
docker run --rm -it -p 8000:8000 -v ${PWD}:/docs:ro,Z -v ${PWD}/mkdocs.yml:/docs/mkdocs.yml:ro,Z --cgroup-manager=systemd squidfunk/mkdocs-material
```
