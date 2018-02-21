# Borrowed from: 
# https://github.com/silven/go-example/blob/master/Makefile
# https://vic.demuzere.be/articles/golang-makefile-crosscompile/

BINARY = addrindex-server
VET_REPORT = vet.report
TEST_REPORT = tests.xml
GOARCH = amd64

VERSION=0.14.1-bitcore
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

# Symlink into GOPATH
GITHUB_USERNAME=jackzampolin
BUILD_DIR=${GOPATH}/src/github.com/${GITHUB_USERNAME}/${BINARY}
CURRENT_DIR=$(shell pwd)
BUILD_DIR_LINK=$(shell readlink ${BUILD_DIR})
ARTIFACT_DIR=build
FLAG_PATH=github.com/${GITHUB_USERNAME}/${BINARY}/cmd
DOCKER_TAG=${VERSION}-$(shell git rev-parse --short HEAD)
DOCKER_IMAGE=quay.io/blockstack/addrindex-server

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-X ${FLAG_PATH}.Version=${VERSION} -X ${FLAG_PATH}.Commit=${COMMIT} -X ${FLAG_PATH}.Branch=${BRANCH}"

# Build the project
all: dep clean linux darwin

linux: dep clean 
	cd ${BUILD_DIR}; \
	GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o ${ARTIFACT_DIR}/${BINARY}-linux-${GOARCH} . ; \
	cd - >/dev/null

darwin:
	cd ${BUILD_DIR}; \
	GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o ${ARTIFACT_DIR}/${BINARY}-darwin-${GOARCH} . ; \
	cd - >/dev/null

dep:
	glide i

docker:
	cd ${BUILD_DIR}; \
	docker build -t ${DOCKER_IMAGE}:${DOCKER_TAG} .
	docker tag ${DOCKER_IMAGE}:${DOCKER_TAG} ${DOCKER_IMAGE}:latest
	docker push ${DOCKER_IMAGE}:${DOCKER_TAG}
	docker push ${DOCKER_IMAGE}:latest

clean:
	-rm -f ${ARTIFACT_DIR}/${BINARY}-*

.PHONY: dep linux darwin fmt clean