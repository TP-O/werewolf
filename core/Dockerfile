FROM golang:1.18.10-alpine3.17

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin ./...

CMD ["uwwolf"]
