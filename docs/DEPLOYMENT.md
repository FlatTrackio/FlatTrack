- [Kubernetes (recommended)](#sec-1)
- [Docker-compose](#sec-2)


# Kubernetes (recommended)<a id="sec-1"></a>

# Docker-compose<a id="sec-2"></a>

The example set up in Docker-compose for FlatTrack is located in [../docker-compose.yml](https://gitlab.com/flattrack/flattrack/-/blob/master/docker-compose.yml). This configuration is not considered production-ready, but for some may be a good enough option. `Important notes`:

-   the docker-compose deployment doesn't include SSL/TLS in the stack, so if you wish to deploy it you will need to add a reverse-proxy to handle SSL/TLS - this isn't something that FlatTrack implements.
-   ensure that the credentials are updated

```sh
docker-compose up -d
```
