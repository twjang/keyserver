BINARY            = keyserver
GITHUB_USERNAME   = twjang
VERSION           = v4.0.2
GOARCH            = amd64
ARTIFACT_DIR      = build
PORT 							= 3000

COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

# Symlink into GOPATH
FLAG_PATH=github.com/${GITHUB_USERNAME}/${BINARY}/cmd

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-X ${FLAG_PATH}.Version=${VERSION} -X ${FLAG_PATH}.Commit=${COMMIT} -X ${FLAG_PATH}.Branch=${BRANCH}"

# Build the project
all: clean linux darwin windows

# Build and Install project into GOPATH using current OS setup
install:
	go install ${LDFLAGS} ./...

test:
	go test -v ./api/...

# Build binary for Linux
linux: clean
	GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o ${ARTIFACT_DIR}/${BINARY}-linux-${GOARCH} . ;

# Build binary for MacOS
darwin:
	GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o ${ARTIFACT_DIR}/${BINARY}-darwin-${GOARCH} . ;

# Build binary for Windows
windows:
	GOOS=windows GOARCH=${GOARCH} go build ${LDFLAGS} -o ${ARTIFACT_DIR}/${BINARY}-windows-${GOARCH}.exe . ;

# Install golang dependencies

# Remove all the built binaries
clean:
	rm -rf ${ARTIFACT_DIR}/*

.PHONY: linux darwin fmt clean
