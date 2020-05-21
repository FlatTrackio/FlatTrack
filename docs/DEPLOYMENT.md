- [Kubernetes (recommended)](#sec-1)
  - [Helm](#sec-1-1)
    - [Configuration](#sec-1-1-1)
    - [Installation](#sec-1-1-2)
- [Docker-compose](#sec-2)


# Kubernetes (recommended)<a id="sec-1"></a>

## Helm<a id="sec-1-1"></a>

### Configuration<a id="sec-1-1-1"></a>

Chart values are located in [../k8s-manifests/helm/values.yaml](https://gitlab.com/flattrack/flattrack/-/blob/master/k8s-manifests/helm/values.yaml).

| Parameter                                             | Description                                                  | Default                                 |
|----------------------------------------------------- |------------------------------------------------------------ |--------------------------------------- |
| replicaCount                                          | FlatTrack replica Pods                                       | 1                                       |
| image.repository                                      | The repo where the image lives                               | registry.gitlab.com/flattrack/flattrack |
| image.tag                                             | Specifies a tag of from the image to use                     | `nil`                                   |
| image.pullPolicy                                      | FlatTrack container pull policy                              | IfNotPresent                            |
| imagePullSecrets                                      | References for the registry secrets to pull FlatTrack from   | `[]`                                    |
| nameOverride                                          | Expand the name of the chart                                 | `""`                                    |
| fullNameOverride                                      | Create a FQDN for the app name                               | `""`                                    |
| podSecurityContext.readOnlyRootFilesystem             | Set the rootfs as read-only                                  | true                                    |
| podSecurityContext.runAsUser                          | The user to run as                                           | 1000                                    |
| podSecurityContext.runAsGroup                         | The group to run as                                          | 1000                                    |
| podSecurityContext.allowPrivilegeEscalation           | If the process in the container can become root              | true                                    |
| service.type                                          | In way which the app is exposed                              | ClusterIP                               |
| service.port                                          | The port to run the app on                                   | 8080                                    |
| horizonalPodAutoscaler.enabled                        | If the Pods should autoscale                                 | true                                    |
| horizonalPodAutoscaler.minReplicas                    | Minimum amount of Pods                                       | 2                                       |
| horizonalPodAutoscaler.maxReplicas                    | Maximum amount of Pods                                       | 5                                       |
| horizonalPodAutoscaler.targetCPUUtilizationPercentage | How much resource should be utilized before scaling          | 5                                       |
| podDisruptionBudget.enabled                           | define the max number of disruptions allowed                 | enable                                  |
| podDisruptionBudget.minAvailable                      | minimum number of Pods that should always be available       | 2                                       |
| podDisruptionBudget.maxUnavailable                    | max amount of Pods that are allowed to be unavailable        |                                         |
| prometheus.enabled                                    | enable Prometheus annotations                                | true                                    |
| extraEnvVars                                          | declare extra environment variables                          |                                         |
| postgres.username                                     | the username for an existing Postgres databse                | flattrack                               |
| postgres.password                                     | the password for an existing Postgres databse                | flattrack                               |
| postgres.host                                         | the host for an existing Postgres databse                    | flattrack                               |
| postgres.database                                     | the database for an existing Postgres databse                | flattrack                               |
| labels                                                | declare labels for all resources                             | `{}`                                    |
| annotations                                           | declare annotations for all resources                        | `{}`                                    |
| ingress.enabled                                       | create an ingress manifests                                  | false                                   |
| ingress.annotations                                   | set annotations for the ingress manifest                     | `{}`                                    |
| ingress.hosts                                         | the hosts which the ingress endpoint should be accessed from |                                         |
| ingress.tls                                           | references to TLS secrets                                    | `[]`                                    |
| resources.limits.cpu                                  | max amount of CPU                                            | 1m                                      |
| resources.limits.memory                               | max amount of memory                                         | 20Mi                                    |
| resources.resources.cpu                               | requested amount of CPU                                      | 1m                                      |
| resources.limits.memory                               | max amount of memory                                         | 20Mi                                    |
| nodeSelector                                          | delcare the node labels for Pod scheduling                   | `{}`                                    |
| tolerations                                           | declare the toleration labels for Pod scheduling             | `[]`                                    |
| affinity                                              | declare the affinity settings for the Pod scheduling         | `{}`                                    |

### Installation<a id="sec-1-1-2"></a>

```sh
# create the namespace
kubectl create namespace flattrack-myflat

# TODO use url instead of path
# install
helm install flattrack-myflat --namespace flattrack-myflat k8s-manifests/helm
```

# Docker-compose<a id="sec-2"></a>

The example set up in Docker-compose for FlatTrack is located in [../docker-compose.yml](https://gitlab.com/flattrack/flattrack/-/blob/master/docker-compose.yml). This configuration is not considered production-ready, but for some may be a good enough option. `Important notes`:

-   the docker-compose deployment doesn't include SSL/TLS in the stack, so if you wish to deploy it you will need to add a reverse-proxy to handle SSL/TLS - this isn't something that FlatTrack implements.
-   ensure that the credentials are updated

```sh
docker-compose up -d
```
