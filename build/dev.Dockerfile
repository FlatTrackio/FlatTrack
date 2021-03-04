FROM alpine:3.11 as extras
RUN apk add tzdata ca-certificates
RUN adduser -D user

FROM scratch
WORKDIR /app
ENV PATH=/app/bin \
  APP_DIST_FOLDER=./dist
COPY web/dist /app/dist
COPY bin/flattrack /app/flattrack
COPY --from=extras /etc/passwd /etc/passwd
COPY --from=extras /etc/group /etc/group
COPY --from=extras /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=extras /etc/ssl /etc/ssl
COPY migrations /app/migrations
COPY templates /app/templates
EXPOSE 8080
USER user
ENTRYPOINT ["/app/bin/flattrack"]
