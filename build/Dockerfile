FROM node:15.9.0-alpine3.13 AS ui
ARG AppBuildVersion="0.0.0"
ARG AppBuildHash="???"
ARG AppBuildDate="???"
ARG AppBuildMode="development"
RUN apk add --no-cache python2 make g++
COPY web /app/web
WORKDIR /app/web
RUN npm i
ENV VUE_APP_AppBuildVersion=$AppBuildVersion \
  VUE_APP_AppBuildHash=$AppBuildHash \
  VUE_APP_AppBuildDate=$AppBuildDate \
  VUE_APP_AppBuildMode=$AppBuildMode
RUN npm run build

FROM golang:1.16.0-alpine3.13 AS api
WORKDIR /app
COPY pkg /app/pkg
COPY cmd /app/cmd
COPY go.* *.go /app/
# TODO build the version, git hash, and build date into the binary from build-args or env
ARG AppBuildVersion="0.0.0"
ARG AppBuildHash="???"
ARG AppBuildDate="???"
ARG AppBuildMode="development"
ARG GOARCH=""
RUN CGO_ENABLED=0 GOOS=linux GOARCH="$GOARCH" go build \
  -a \
  -installsuffix cgo \
  -ldflags "-extldflags '-static' -s -w \
    -X gitlab.com/flattrack/flattrack/pkg/common.AppBuildVersion=$AppBuildVersion \
    -X gitlab.com/flattrack/flattrack/pkg/common.AppBuildHash=$AppBuildHash \
    -X gitlab.com/flattrack/flattrack/pkg/common.AppBuildDate=$AppBuildDate \
    -X gitlab.com/flattrack/flattrack/pkg/common.AppBuildMode=$AppBuildMode" \
  -o bin/flattrack \
  main.go

FROM alpine:3.11 as extras
RUN apk add --no-cache tzdata ca-certificates
RUN adduser -D user

FROM scratch
LABEL maintainer="Caleb Woodbine <calebwoodbine.public@gmail.com>" \
  org.opencontainers.image.title="flattrack" \
  org.opencontainers.image.description="Collaborate with your flatmates" \
  org.opencontainers.image.url="https://gitlab.com/flattrack/flattrack" \
  org.opencontainers.image.source="git@gitlab.com/flattrack/flattrack" \
  org.opencontainers.image.vendor="FlatTrack" \
  org.opencontainers.image.licenses="AGPL-3.0" \
  org.opencontainers.image.documentation="https://flattrack.gitlab.io/flattrack" \
  org.opencontainers.image.authors="FlatTrack authors"
WORKDIR /app
ENV PATH=/app \
  APP_DIST_FOLDER=./dist
COPY --from=ui /app/web/dist /app/dist
COPY --from=api /app/bin/flattrack /app/bin/flattrack
COPY --from=extras /etc/passwd /etc/passwd
COPY --from=extras /etc/group /etc/group
COPY --from=extras /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=extras /etc/ssl /etc/ssl
COPY migrations /app/migrations
COPY templates /app/templates
EXPOSE 8080
USER user
ENTRYPOINT ["/app/bin/flattrack"]
