FROM node:19.5.0 as frontend

WORKDIR /app/src

COPY package*.json ./
RUN npm install

COPY ./ ./
RUN npm run build

FROM nginx:1.18.0
COPY --from=frontend /app/src/dist/client/ /usr/share/nginx/html
