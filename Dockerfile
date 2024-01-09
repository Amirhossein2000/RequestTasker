FROM golang:latest

WORKDIR /app

COPY . .

RUN mkdir -p ./bin/ && go build -o ./bin ./cmd/

EXPOSE 8080

CMD ["./bin/cmd"]
