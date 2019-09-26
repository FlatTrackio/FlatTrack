FROM node:12.10.0
RUN mkdir -p /opt/flattrack
ADD . /opt/flattrack/.
WORKDIR /opt/flattrack
ENV NODE_ENV=production
ENV APP_PORT=80
CMD npm start
