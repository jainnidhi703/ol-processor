FROM golang:latest

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o ./cmd

EXPOSE 3000

CMD ["./app/main"]