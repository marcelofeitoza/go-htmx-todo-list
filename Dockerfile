FROM golang:1.21

WORKDIR /app

COPY . /app

RUN CGO_ENABLED=0 GOOS=linux go mod download && go build -o main . && chmod +x main

EXPOSE 8080

CMD ["./main"]
