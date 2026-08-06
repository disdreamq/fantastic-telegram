FROM golang:1.26.5-alpine AS builder

WORKDIR /app

COPY ./services/notification/go.mod ./services/notification/go.sum ./
RUN go mod download

COPY ./services/notification .


RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /main ./cmd

FROM alpine:latest

RUN adduser -D -g '' appuser
USER appuser

WORKDIR /app

COPY --from=builder /main .

EXPOSE 8082

CMD ["./main"]