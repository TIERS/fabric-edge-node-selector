FROM golang:1.17.2-alpine as development

WORKDIR /app

COPY go.mod go.sum .env ./
RUN go mod download

# For DEBUGGING purposes
# RUN apk update && apk add curl

COPY . .

RUN go install github.com/cespare/reflex@latest

EXPOSE 8080

CMD reflex -g '*.go' go run cmd/main.go --start-service