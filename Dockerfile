# Base Stage
FROM golang:1.22-alpine AS base
WORKDIR /app
COPY ./go.mod ./go.sum ./
RUN mkdir -p dist && \
    go mod download

# Debug in Docker (only for backend)
FROM alpine:latest as debug
WORKDIR /
COPY . .
COPY --from=base /usr/local/go/ /usr/local/go/
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOBIN=/usr/local/go/bin/
RUN go install github.com/go-delve/delve/cmd/dlv@latest &&  \
    go build -gcflags="all=-N -l" -o dist/bin cmd/main.go && \
    apk add --no-cache bash curl libxml2-utils jq
CMD ["/usr/local/go/bin/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "dist/bin"]

# Build Production Stage
FROM base as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o dist/bin cmd/main.go

# Build Web
FROM node:20-alpine as web
RUN mkdir -p /web/build
WORKDIR /web
COPY ./web .
RUN yarn install && yarn build

# Production Stage
FROM alpine:latest as prod
WORKDIR /
COPY --from=build /app/dist/bin ./bin
COPY --from=web /web/build ./web/build
RUN apk add --no-cache bash curl libxml2-utils jq
EXPOSE 8080 8088
CMD ["./bin"]