FROM node:12.10.0-alpine AS ui
WORKDIR /app
COPY src /app/src
COPY public /app/public
COPY *.js *.json /app/
RUN npm i
RUN npm run build:frontend

FROM golang:1.13.4-alpine3.10 AS api
WORKDIR /app
COPY src /app/src
COPY go.* /app/
RUN adduser -D user
# TODO build the version, git hash, and build date into the binary from build-args or env
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static" -s -w' -o flattrack src/backend/main.go

FROM alpine:3.10
WORKDIR /app
ENV PATH=/app
COPY --from=ui /app/dist /app/dist
COPY --from=ui /app/package.json .
COPY --from=api /app/flattrack .
COPY --from=api /etc/passwd /etc/passwd
COPY --from=api /etc/group /etc/group
COPY --chown=user migrations /app/migrations
EXPOSE 8080
USER user
ENTRYPOINT ["/app/flattrack"]
