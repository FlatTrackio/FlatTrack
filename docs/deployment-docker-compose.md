# Deployment with Docker Compose

The example set up in Docker-compose for FlatTrack is located in [../deployments/docker-compose.yml](https://gitlab.com/flattrack/flattrack/-/blob/master/deployments/docker-compose.yml). This configuration is not considered production-ready, but for some may be a good enough option.

Deploy the docker-compose stack:

    docker-compose up -d

`Important notes`:

-   the docker-compose deployment doesn&rsquo;t include SSL/TLS in the stack, so if you wish to deploy it you will need to add a reverse-proxy to handle SSL/TLS - this isn&rsquo;t something that FlatTrack implements.
-   ensure that the credentials are updated

# Extra notes

To configure FlatTrack, please refer to the [configuration guide](./configuration.md).

