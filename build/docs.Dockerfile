FROM theiceshelf/firn AS build
WORKDIR /app/site
COPY --chown=user:user docs /app/site
RUN firn build

FROM registry.gitlab.com/safesurfer/go-http-server:1.0.1
COPY --from=build /app/site/_firn/_site /app/site
