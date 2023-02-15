FROM node:19.5.0 as frontend

WORKDIR /app/src

ARG MYSQL_USER
ARG MYSQL_PASSWORD
ARG MYSQL_DB_NAME
ARG MYSQL_TEST_DB_NAME
ARG MYSQL_HOST
ARG API_PORT
ARG API_HOST

COPY package*.json ./
RUN npm install

COPY ./ ./
RUN npm run build

FROM nginx:1.18.0
COPY --from=frontend /app/src/dist/client/ /usr/share/nginx/html
