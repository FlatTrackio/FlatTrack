FROM node:12.13.1-alpine3.10 AS fetchdeps
WORKDIR /app
ADD . .
RUN npm i
RUN npm rebuild node-sass
RUN npm run build
RUN adduser -D flattrack

FROM gcr.io/distroless/nodejs
ENV PUID=1000 \
    PGID=1000 \
    NODE_ENV=production \
    APP_PORT=8080 \
    NO_UPDATE_NOTIFIER=true
COPY --from=fetchdeps /app /app
COPY --from=fetchdeps /etc/passwd /etc/passwd
WORKDIR /app
USER flattrack
CMD ["src/server/server.js"]
