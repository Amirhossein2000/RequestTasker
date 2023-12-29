FROM golang:latest

WORKDIR /app

COPY . .

RUN go build -o ./bin ./cmd/

EXPOSE 8080

CMD ["./bin/cmd"]
