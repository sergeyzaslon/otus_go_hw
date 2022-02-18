FROM golang:1.16.2 as build

ENV CGO_ENABLED=0
RUN go install github.com/pressly/goose/v3/cmd/goose@v3.5.0
RUN go install github.com/antelman107/net-wait-go@latest

FROM alpine:3.9
RUN apk add --update --no-cache netcat-openbsd

COPY --from=build /go/bin/goose /bin/goose
COPY --from=build /go/bin/net-wait-go /bin/net-wait-go

WORKDIR /var/migrations

COPY migrations/* /var/migrations/

ENV DRIVER=postgres
ENV DSN=""

CMD ["/bin/goose", "up"]