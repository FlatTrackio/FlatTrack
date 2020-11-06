FROM node:14.7.0-alpine3.11 AS ui
LABEL maintainer="Caleb Woodbine <calebwoodbine.public@gmail.com>"
ARG AppBuildVersion="0.0.0"
ARG AppBuildHash="???"
ARG AppBuildDate="???"
ARG AppBuildMode="development"
RUN apk add python2 make g++
COPY web /app/web
WORKDIR /app/web
RUN npm i
ENV VUE_APP_AppBuildVersion=$AppBuildVersion \
  VUE_APP_AppBuildHash=$AppBuildHash \
  VUE_APP_AppBuildDate=$AppBuildDate \
  VUE_APP_AppBuildMode=$AppBuildMode
RUN npm run build

FROM golang:1.14.6-alpine3.11 AS api
WORKDIR /app
COPY src /app/src
COPY go.* /app/
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
    -X gitlab.com/flattrack/flattrack/src/backend/common.AppBuildVersion=$AppBuildVersion \
    -X gitlab.com/flattrack/flattrack/src/backend/common.AppBuildHash=$AppBuildHash \
    -X gitlab.com/flattrack/flattrack/src/backend/common.AppBuildDate=$AppBuildDate \
    -X gitlab.com/flattrack/flattrack/src/backend/common.AppBuildMode=$AppBuildMode" \
  -o flattrack \
  src/backend/main.go

FROM alpine:3.11 as extras
RUN apk add tzdata ca-certificates
RUN adduser -D user

FROM scratch
WORKDIR /app
ENV PATH=/app
COPY --from=ui /app/web/dist /app/dist
COPY --from=ui /app/web/package.json .
COPY --from=api /app/flattrack .
COPY --from=extras /etc/passwd /etc/passwd
COPY --from=extras /etc/group /etc/group
COPY --from=extras /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=extras /etc/ssl /etc/ssl
COPY migrations /app/migrations
COPY templates /app/templates
EXPOSE 8080
USER user
ENTRYPOINT ["/app/flattrack"]
