# alipine → debian-slimの予定

FROM golang:1.23.4-alpine

WORKDIR /usr/src

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

CMD ["air", "-c", ".air.toml"]
