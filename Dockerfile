# syntax=docker/dockerfile:1

FROM golang:1.22.4

WORKDIR /usr/src/msgsrv

COPY **/go.mod **/go.sum internal/**/go.mod internal/**/go.sum ./
RUN go mod download && go mod verify

COPY . .

WORKDIR /usr/src/msgsrv/cmd/msgsrv-redis

RUN go build -o /usr/local/bin/msgsrv-redis

EXPOSE 80

CMD ["msgsrv-redis"]