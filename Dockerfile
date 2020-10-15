FROM golang:1.11-alpine AS build

# Install tools required for project
# Run `docker build --no-cache .` to update dependencies
RUN apk add --no-cache git gcc libc-dev

# List project dependencies with Gopkg.toml and Gopkg.lock
# These layers are only re-built when Gopkg files are updated
COPY . /app
WORKDIR /app
RUN GOOS=linux CGO_ENABLED=0 go build -o /app/build/hyper-cas main.go

# This results in a single layer image
FROM alpine:3.7
COPY ./docker-hyper-cas.yaml /app/hyper-cas.yaml
COPY --from=build /app/build/hyper-cas /app/hyper-cas
WORKDIR /app
ENTRYPOINT ["/app/hyper-cas"]
CMD ["--config", "/app/hyper-cas.yaml", "serve"]
