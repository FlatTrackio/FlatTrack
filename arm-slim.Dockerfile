FROM arm64v8/node:12.7.0-alpine
RUN mkdir -p /opt/flattrack
ADD . /opt/flattrack/.
WORKDIR /opt/flattrack
ENV NODE_ENV=production
ENV APP_PORT=80
CMD npm start
