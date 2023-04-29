FROM golang:alpine

WORKDIR /app

COPY . .

RUN go mod download && go mod tidy

RUN go build -o productApp .

EXPOSE 8083

CMD ["./productApp"]
