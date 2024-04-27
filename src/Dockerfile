# Description: Dockerfile for building the Go and next JS application

# Go build stage
FROM golang:1.22.2 AS go-builder

WORKDIR /be/go

COPY be/Go .

RUN go build -o app

# Set CGO_ENABLED=0 to disable cgo which gives us a static binary.
RUN CGO_ENABLED=0 go build -o app

# Next JS build stage
FROM node:20.11.0 AS js-builder

WORKDIR /fe

COPY fe/package.json .
COPY fe .

RUN npm install
RUN npm run build

# Final stage
FROM ubuntu:20.04

RUN apt-get update && apt-get install -y supervisor curl
RUN curl -sL https://deb.nodesource.com/setup_18.x | bash -
RUN apt-get install -y nodejs

COPY --from=go-builder /be/go/app /app/be/go/app
COPY --from=js-builder /fe /app/fe

COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

CMD ["/usr/bin/supervisord"]