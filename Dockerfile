FROM node:12.10.0
RUN mkdir -p /opt/flattrack
ADD . /opt/flattrack/.
RUN chown -R node:node /opt/flattrack
USER node
WORKDIR /opt/flattrack
RUN npm i
RUN npm run build
ENV PUID=1000 \
    PGID=1000 \
    NODE_ENV=production \
    APP_PORT=8080 \
    NO_UPDATE_NOTIFIER=true
CMD npm start
