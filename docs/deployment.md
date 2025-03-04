# Deployment

supported deployment methods:

-   [Kubernetes](./deployment-kubernetes.md)
-   [Docker-Compose](./deployment-docker-compose.md)

## Verifying container images

Before deployment, verify the FlatTrack container image with the following

```shell
cosign verify \
  --certificate-oidc-issuer 'https://gitlab.com' \
  --certificate-identity-regexp 'https://gitlab.com/flattrack/flattrack//.gitlab-ci.yml@refs/(heads/main|tags/.*)' \
  -o text \
  registry.gitlab.com/flattrack/flattrack:latest
```

note that the image can also be replaced with an image hosted on __docker.io/flattrack/flattrack_.

(read about: [cosign and sigstore docs](https://docs.sigstore.dev/about/overview/))

## Fetch SBOM

To verify SBOM and get it's contents for an image on a given architecture (such as linux/arm64), use a command like the following:

```shell
IMAGE=registry.gitlab.com/flattrack/flattrack
TAG=latest
PLATFORM=linux/arm64
DIGEST="$(crane digest --platform "$PLATFORM" "$IMAGE:$TAG")"
cosign \
    verify-attestation \
    "$IMAGE@$DIGEST" \
    --certificate-identity-regexp 'https://gitlab.com/flattrack/flattrack//.gitlab-ci.yml@refs/(heads/main|tags/.*)' \
    --certificate-oidc-issuer 'https://gitlab.com' \
      | jq -r .payload | base64 -D | jq . | jq -r .predicate.Data
```

Changing TAG and PLATFORM to the supported & desired values.
