# build stage
FROM golang:1-alpine AS build-env

# Install make
RUN apk add --update make git

# Install Glide
RUN go get -u github.com/Masterminds/glide/...

# Create build dir
RUN mkdir -p /go/src/github.com/jackzampolin/addrindex-server

# Work out of build dir
WORKDIR /go/src/github.com/jackzampolin/addrindex-server

# Copy in source
COPY . .

# Build app
RUN make linux


# Production Image
FROM alpine

COPY --from=build-env /go/src/github.com/jackzampolin/addrindex-server/build/addrindex-server-linux-amd64 /usr/bin/addrindex-server

COPY config.sample.yaml /root/.addrindex-server.yaml

ENTRYPOINT ["/usr/bin/addrindex-server"]

CMD ["serve"]
