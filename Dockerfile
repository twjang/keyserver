FROM golang:alpine AS build-env

# Injest build args from Makefile
ARG BINARY
ARG GITHUB_USERNAME
ARG GOARCH
ENV BINARY=${BINARY}

# Set up dependencies
ENV PACKAGES make git curl

# Set working directory for the build
WORKDIR /go/src/github.com/${GITHUB_USERNAME}/${BINARY}

# Install dependencies
RUN apk add --update $PACKAGES

# Install dep
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# Force the go compiler to use modules 
ENV GO111MODULE=on

# Add source files
COPY . .

# Make the binary
RUN make install

# Final image
FROM alpine:edge

# Install ca-certificates
RUN apk add --update ca-certificates
WORKDIR /root

# Copy over binaries from the build-env
COPY --from=build-env /go/bin/${BINARY} /usr/bin/${BINARY}

# Run ${BINARY} by default
CMD ${BINARY}
