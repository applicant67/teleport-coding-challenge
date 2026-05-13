FROM node:12

WORKDIR /app

COPY ./frontend /app
COPY ./certs /certs

RUN npm install

ENTRYPOINT npm run start