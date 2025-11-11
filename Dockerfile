FROM golang:1.25 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o userapp

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/userapp /app/userapp

EXPOSE 8080

CMD ["./userapp"]