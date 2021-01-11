FROM golang:1.15-alpine AS build
LABEL maintainer="vinayakshnd@gmail.com"

WORKDIR /go/src/github.com/vinayakshnd/gogetgithub

# Install setup dependencies
RUN apk update && \
    apk add git && \
    apk add make

COPY main.go main.go
COPY httphandlers httphandlers
COPY utils utils
COPY Makefile Makefile

COPY .git .git

# Compile
RUN mkdir bin && \
    make deps && \
    make build

# Build a fresh container with just the binaries
FROM alpine

COPY --from=build /go/src/github.com/vinayakshnd/gogetgithub bin
ENTRYPOINT ["gogetgithub"]