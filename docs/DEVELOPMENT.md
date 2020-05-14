- [Development cycle](#sec-1)
  - [Bring up the deployments](#sec-1-1)
  - [Port-forward the database connection](#sec-1-2)
  - [Quick development patterns](#sec-1-3)
    - [Backend](#sec-1-3-1)
    - [Frontend](#sec-1-3-2)
  - [Auto build and launch with [Tilt](https://tilt.dev)](#sec-1-4)
- [Additional](#sec-2)
  - [Manually connecting to the Postgres database](#sec-2-1)
  - [Remove migrations](#sec-2-2)


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
npm run build:frontend
```

## Auto build and launch with [Tilt](https://tilt.dev)<a id="sec-1-4"></a>

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
