FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /ratings-service .

EXPOSE 8003
CMD ["/ratings-service"]
