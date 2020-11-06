FROM squidfunk/mkdocs-material AS build
WORKDIR /app
ADD docs docs
ADD mkdocs.yml /app/
RUN mkdocs build

FROM registry.gitlab.com/safesurfer/go-http-server:1.0.0
COPY --from=build /app/site /app/site
