# Deployment with Kubernetes

## Configuration

To configure FlatTrack environment variables, please refer to the [configuration guide](./configuration.md).

## Installation

First, create a Postgres database.

Create the Namespace:

    kubectl create namespace flattrack-myflat

Define the Deployment:

    ---
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: flattrack
      namespace: flattrack-myflat
    spec:
      selector:
        matchLabels:
          app: flattrack
      replicas: 1
      template:
        metadata:
          labels:
            app: flattrack
        spec:
          containers:
            - name: flattrack
              image: registry.gitlab.com/flattrack/flattrack:latest # TODO use release tag
              securityContext:
                allowPrivilegeEscalation: false
                capabilities:
                  drop:
                    - ALL
                readOnlyRootFilesystem: true
                runAsGroup: 1000
                runAsUser: 1000
                seccompProfile:
                  type: RuntimeDefault
              ports:
                - containerPort: 8080
              env:
                - name: APP_DB_HOST
                  value: postgres
                - name: APP_DB_USERNAME
                  value: postgres
                - name: APP_DB_PASSWORD
                  value: postgres
                - name: APP_DB_DATABASE
                  value: postgres
                - name: APP_DB_SSLMODE
                  value: disable
                  
Next steps:

- create a Service
- create an Ingress, VirtualService or HTTPRoute
- setup TLS, such as including LetsEncrypt

## Installation with Knative

First, create a Postgres database.

Create the Service:

    kn -n flattrack-myflat service create flattrack \
      --env APP_DB_USER=postgres \
      --env APP_DB_PASSWORD=postgres \
      --env APP_DB_DATABASE=postgres \
      --env APP_DB_HOST=postgres \
