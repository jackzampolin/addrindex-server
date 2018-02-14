# Borrowed from: 
# https://github.com/silven/go-example/blob/master/Makefile
# https://vic.demuzere.be/articles/golang-makefile-crosscompile/

BINARY = addrindex-server
VET_REPORT = vet.report
TEST_REPORT = tests.xml
GOARCH = amd64

VERSION=0.13.0.2-counterparty
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

# Symlink into GOPATH
GITHUB_USERNAME=jackzampolin
BUILD_DIR=${GOPATH}/src/github.com/${GITHUB_USERNAME}/${BINARY}
CURRENT_DIR=$(shell pwd)
BUILD_DIR_LINK=$(shell readlink ${BUILD_DIR})
ARTIFACT_DIR=build
FLAG_PATH=github.com/${GITHUB_USERNAME}/${BINARY}/cmd

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-X ${FLAG_PATH}.Version=${VERSION} -X ${FLAG_PATH}.Commit=${COMMIT} -X ${FLAG_PATH}.Branch=${BRANCH}"

# Build the project
all: dep clean linux darwin

linux: 
	cd ${BUILD_DIR}; \
	GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o ${ARTIFACT_DIR}/${BINARY}-linux-${GOARCH} . ; \
	cd - >/dev/null

darwin:
	cd ${BUILD_DIR}; \
	GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o ${ARTIFACT_DIR}/${BINARY}-darwin-${GOARCH} . ; \
	cd - >/dev/null

dep:
	glide i

clean:
	-rm -f ${ARTIFACT_DIR}/${BINARY}-*

.PHONY: dep linux darwin fmt clean