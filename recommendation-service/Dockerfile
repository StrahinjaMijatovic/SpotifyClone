FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /recommedation-service .

EXPOSE 8006
CMD ["/recommedation-service"]
