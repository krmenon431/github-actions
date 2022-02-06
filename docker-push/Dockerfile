# argument for Go version
ARG GO_VERSION=1.17.6

# STAGE 1: building the executable
FROM golang:1.17.6-alpine3.15 as build

# git required for go mod
RUN apk add --no-cache git ca-certificates
RUN apk --update --no-cache upgrade && apk add --no-cache ca-certificates

# Working directory will be created if it does not exist
WORKDIR /src

# We use go modules; copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go help


# Import code
COPY . ./

# # Run tests
# RUN CGO_ENABLED=0 go test -timeout 30s -v github.com/gbaeke/go-template/pkg/api

# Build the executable
RUN CGO_ENABLED=0 go build

# STAGE 2: build the container to run
FROM gcr.io/distroless/static AS final

# add maintainer label
LABEL maintainer="krishnanunni"

# copy compiled app
COPY --from=build /src/docker-push .

# run binary
ENTRYPOINT ["./docker-push"]