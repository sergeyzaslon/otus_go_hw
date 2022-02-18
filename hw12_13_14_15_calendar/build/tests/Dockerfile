# Собираем в гошке
FROM golang:1.16.2 as build

RUN go install github.com/antelman107/net-wait-go@latest

WORKDIR "/opt/calendar"

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . "/opt/calendar/"

# CMD ["go", "test", "github.com/sergeyzaslon/otus_go_hw/hw12_13_14_15_calendar/cmd/calendar_tests"]
