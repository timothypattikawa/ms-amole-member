FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN ENV=staging go build -o main cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY application-*.yml .

EXPOSE 9091 9091
CMD [ "/app/main" ]
