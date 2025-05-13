FROM golang:bookworm AS build

WORKDIR /usr/src/fileman-net

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /bin/fileman-net-client ./cmd/client

FROM debian:bookworm AS result

COPY --from=build /bin/fileman-net-client /bin/fileman-net-client

ENTRYPOINT [ "/bin/fileman-net-client", "--port", "12334" ]
CMD [ "--address", "127.0.0.1" ]
