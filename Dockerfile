FROM golang:latest as builder

WORKDIR /app
ADD . /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main

FROM alpine:latest AS production

ENV TZ=Europe/Moscow
RUN apk update
RUN apk upgrade
RUN apk add ca-certificates && update-ca-certificates
RUN apk add --update tzdata
RUN rm -rf /var/cache/apk/*

ENV APPLICATION=production
WORKDIR /app
COPY --from=builder /app .
EXPOSE 8080
CMD ["./main"]