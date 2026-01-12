FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /subsriptions-service .

EXPOSE 8004
CMD ["/subsriptions-service"]
