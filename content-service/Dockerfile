FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /content-service .

EXPOSE 8002
CMD ["/content-service"]
