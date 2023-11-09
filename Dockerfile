FROM golang:1.21

WORKDIR /app

RUN apt-get update && apt-get install graphviz -y

COPY go.mod  go.sum ./
RUN go mod download

COPY . .

RUN go build cmd/ol-processor/main.go

EXPOSE 3000

CMD ["/app/main"]