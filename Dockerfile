FROM golang:1.21.5-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /smart-energy-price-update

ENV CRON_SCHEDULE_PRICE "5 * * * *"

ENV CRON_SCHEDULE_CONSUMPTION "0 6 * * *"

CMD ["/smart-energy-price-update"]