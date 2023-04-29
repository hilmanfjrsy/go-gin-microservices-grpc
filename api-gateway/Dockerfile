FROM golang:alpine

WORKDIR /app

COPY . .

RUN go mod download && go mod tidy

RUN go build -o apiGateway .

EXPOSE 8081

CMD ["./apiGateway"]
