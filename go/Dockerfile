FROM golang:tip-alpine3.23 AS builder

#create a working directory/folder inside the container
WORKDIR /app

# Copy Go modules and dependencies
COPY go.mod go.sum ./
RUN go mod download

EXPOSE 8080

#copy the source code
# COPY . .

# RUN go build -o app ./cmd/server/

# CMD [ "./app" ]
