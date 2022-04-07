FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

EXPOSE 5000

COPY . .

RUN ["go", "build", "cmd/main.go"]

CMD ["./main"]