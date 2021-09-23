FROM golang:alpine AS builder

# Set necessary environment variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOROOT=/usr/local/go \
    GOPATH=/go

ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH

RUN go env -w GOPRIVATE=github.com/Beep-Technologies

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .

RUN apk add git

ARG GITHUB_TOKEN

RUN git config --global url."https://b33pb0t:$GITHUB_TOKEN@github.com".insteadOf "https://github.com"

RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o main ./cmd/bb3-ocpp-ws

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy build to main folder
RUN cp /build/main .

# Build a small image
FROM alpine:3.12

RUN apk update && apk add ca-certificates \
    && apk add --update bash gzip \
    && update-ca-certificates \
    && rm -rf /var/cache/apk/*

ENV LANG=C.UTF-8

COPY --from=builder /dist/main /
 
# Command to run
ENTRYPOINT ["/main"]
