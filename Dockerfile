FROM golang:latest as builder

WORKDIR /app
ADD . /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main

FROM alpine:latest AS production
WORKDIR /app
COPY --from=builder /app .
EXPOSE 8080
CMD ["./main"]