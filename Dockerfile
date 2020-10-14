FROM golang:1.11-alpine AS build

# Install tools required for project
# Run `docker build --no-cache .` to update dependencies
RUN apk add --no-cache git gcc libc-dev

# List project dependencies with Gopkg.toml and Gopkg.lock
# These layers are only re-built when Gopkg files are updated
COPY . /app
WORKDIR /app
RUN go build -o /app/build/hyper-cas

# This results in a single layer image
FROM scratch
COPY docker-hyper-cas.yaml /app/hyper-cas.yaml
COPY --from=build /app/build/hyper-cas /app/hyper-cas
WORKDIR /app
ENTRYPOINT ["hyper-cas"]
CMD [" serve --config hyper-cas.yaml"]
