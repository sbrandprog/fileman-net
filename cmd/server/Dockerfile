FROM golang:bookworm AS build

WORKDIR /usr/src/fileman-net

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /bin/fileman-net-server ./cmd/server

FROM debian:bookworm AS result

COPY --from=build /bin/fileman-net-server /bin/fileman-net-server

EXPOSE 12334/tcp
VOLUME [ "/data" ]

ENTRYPOINT [ "/bin/fileman-net-server", "--port", "12334", "--server-wd", "/data" ]
