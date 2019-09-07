FROM node:12.10.0-alpine
RUN mkdir -p /opt/flattrack
ADD . /opt/flattrack/.
WORKDIR /opt/flattrack
ENV NODE_ENV=production
CMD npm start