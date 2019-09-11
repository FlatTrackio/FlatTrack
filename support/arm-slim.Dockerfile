FROM arm64v8/node:12.7.0-alpine
RUN mkdir -p /opt/flattrack
ADD . /opt/flattrack/.
WORKDIR /opt/flattrack
ENV NODE_ENV=production
CMD npm start