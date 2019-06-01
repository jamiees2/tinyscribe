# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from golang v1.11 base image
FROM golang:1.12 as builder

# Add Maintainer Info
LABEL maintainer="James Sigurdarson <jamiees2@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /go/src/github.com/jamiees2/tinyscribe

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

# Download dependencies
RUN go get -d -v ./...

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/tinyscribe .


######## Start a new stage from scratch #######
FROM alpine:latest  

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /go/bin/tinyscribe .

ENV LISTEN_HOST ":1463"

CMD ["./tinyscribe"] 
