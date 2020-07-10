# build stage
FROM golang as builder

ENV GO111MODULE=on

ADD . /app
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build


EXPOSE 8080
ENTRYPOINT ["/app/fReader"]