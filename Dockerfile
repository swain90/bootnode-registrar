# Use Go 1.20
FROM golang:1.20-alpine AS builder

ENV GOPATH=/go
ENV APPROOT=${GOPATH}/src/github.com/jpoon/bootnode-registrar
ENV GOARCH=amd64  
# Set the target architecture
ENV GOOS=linux    
# Set the target operating system

# Package dependencies
RUN apk add --update --no-cache git gcc libc-dev ca-certificates file

# Enable Go modules
ENV GO111MODULE=on

# Set Go proxy
ENV GOPROXY=https://proxy.golang.org,direct

# Copy go.mod and go.sum files first
WORKDIR ${APPROOT}
COPY go.mod go.sum ./

# Install and update CA certificates before downloading dependencies
RUN apk add --no-cache ca-certificates && update-ca-certificates

# Download dependencies
RUN go mod download

# Copy the rest of the project files
COPY . .

# Ensure go.mod is in the correct location
RUN ls -la ${APPROOT}

# Get dependencies and compile
RUN go mod tidy -compat=1.17
RUN GOOS=${GOOS} GOARCH=${GOARCH} go build -o bootnode-registrar

# Verify the binary format
RUN file bootnode-registrar

# Use a minimal base image for the final stage
FROM alpine:latest

# Install CA certificates in the final image as well
RUN apk add --no-cache ca-certificates file

WORKDIR /work
COPY --from=builder /go/src/github.com/jpoon/bootnode-registrar/bootnode-registrar /work/bootnode-registrar

# Verify the binary format in the final image
RUN file /work/bootnode-registrar

ENTRYPOINT [ "./bootnode-registrar" ]
EXPOSE 9898