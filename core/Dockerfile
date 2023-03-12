FROM golang:1.18.10-alpine3.17 as build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /app/uwwolf

FROM alpine:3.16.4

COPY --from=build-stage /app/uwwolf /usr/local/bin

CMD ["uwwolf"]
