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
FROM golang:1.17-alpine AS go
WORKDIR /build

RUN apk add --no-cache build-base

COPY go.mod go.sum ./
RUN go mod download

RUN go build github.com/mattn/go-sqlite3

COPY . .
COPY ./internal/server/http/frontend/templates ./internal/server/http/frontend/templates 
COPY --from=node /frontend/public ./internal/server/http/frontend/public

RUN go build -ldflags '-extldflags "-static"' -o ./epigram-server ./cmd/server

# Distroless final container
FROM gcr.io/distroless/base
WORKDIR /server

COPY --from=go /build/epigram-server .

EXPOSE 8080

USER nonroot:nonroot

CMD [ "./epigram-server" ]