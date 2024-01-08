# Start from the latest golang base image
FROM golang:latest

# Add Maintainer Info
LABEL maintainer="AMIT AMRUTIYA <amitamrutiya2210@gmail.com>"

RUN mkdir /app

ADD . /app

# Set the Current Working Directory inside the container
WORKDIR /app

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go build -o main .

# Command to run the executable
CMD ["/app/main"]