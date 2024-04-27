# Next JS build stage
FROM node:20.11.0 AS js-builder

WORKDIR .

COPY . .

RUN npm install
RUN npm run build