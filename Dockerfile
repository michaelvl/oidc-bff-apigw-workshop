FROM node:16.10.0-alpine3.12

RUN mkdir -p /apps/api-gw /apps/protected-api /apps/cdn /apps/spa

COPY api-gw /apps/api-gw
RUN cd /apps/api-gw && npm install

COPY cdn /apps/cdn
RUN cd /apps/cdn && npm install

COPY protected-api /apps/protected-api
RUN cd /apps/protected-api && npm install

COPY spa/dist /apps/spa

EXPOSE 5000

WORKDIR /apps/cdn
CMD [ "npm", "start" ]
