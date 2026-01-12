FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /notifications-service .

EXPOSE 8005
CMD ["/notifications-service"]
