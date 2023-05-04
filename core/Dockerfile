FROM golang:1.18.10-alpine3.17 as build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /app/uwwolf

FROM alpine:3.16.4

LABEL maintainer="Le Tran Phong <letranphong2k1@gmail.com>"

COPY --from=build-stage /app/uwwolf /usr/local/bin

RUN addgroup -S ww -g 1000 && adduser -S ww -u 1000 -G ww
RUN chown ww:ww /usr/local/bin/uwwolf

USER ww

CMD ["uwwolf"]
