# syntax=docker/dockerfile:1
# Epigram multi-stage docker build process

# Node - build tailwind css files
FROM node:16-slim AS node
WORKDIR /frontend

COPY ./internal/server/http/frontend/package*.json ./
RUN npm install

COPY ./internal/server/http/frontend ./
RUN npm run build --production

# Go - compile go project
FROM golang:1.20-alpine AS go
WORKDIR /build

RUN apk add --no-cache build-base

COPY go.mod go.sum ./
RUN go mod download

RUN go build github.com/mattn/go-sqlite3

COPY ./cmd/server ./cmd/server
COPY ./internal ./internal
COPY --from=node /frontend/public ./internal/server/http/frontend/public

RUN go build -ldflags '-extldflags "-static"' -o ./epigram-server ./cmd/server

# Distroless final container
FROM gcr.io/distroless/base
WORKDIR /server

COPY --from=go /build/epigram-server .

ENV EP_PORT=80
ENV EP_ADDRESS=0.0.0.0
EXPOSE 80

USER nonroot:nonroot

CMD [ "./epigram-server" ]