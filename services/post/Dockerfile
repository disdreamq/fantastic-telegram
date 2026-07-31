FROM golang:1.26.5-alpine AS builder

WORKDIR /app

COPY ./services/post/go.mod ./services/post/go.sum ./
RUN go mod download

COPY ./services/post .


RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /main ./cmd

FROM alpine:latest

RUN adduser -D -g '' appuser
USER appuser

WORKDIR /app

COPY --from=builder /main .

EXPOSE 8080

CMD ["./main"]