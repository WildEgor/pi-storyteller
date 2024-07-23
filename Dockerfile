# Base Stage
FROM golang:1.22-alpine AS base
WORKDIR /app
COPY ./go.mod ./go.sum ./
RUN mkdir -p dist && \
    go mod download

# Development Stage
FROM base as dev
WORKDIR /app/
COPY . .
RUN go install github.com/air-verse/air@latest && \
    go build -o dist/app cmd/main.go
CMD ["air", "-c", ".air-unix.toml", "-d"]

# Debug in Docker (only for backend)
FROM base as debug
WORKDIR /
COPY . .
RUN go install github.com/go-delve/delve/cmd/dlv@latest &&  \
    go build -gcflags="all=-N -l" -o /app cmd/main.go &&  \
    mv /go/bin/dlv /
CMD ["/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/app"]

# Build Production Stage
FROM base as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o dist/app cmd/main.go

# Build Web
FROM node:20-alpine as web
WORKDIR /web/
COPY ./web .
RUN yarn install && yarn build

# Production Stage
FROM alpine:latest as prod
WORKDIR /
COPY --from=build /app/dist/app ./app
COPY --from=web /web/build ./web/build
EXPOSE 8080 8088
CMD ["./app"]