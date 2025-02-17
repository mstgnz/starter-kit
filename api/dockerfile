FROM golang:1.23-alpine as builder
RUN apk add --no-cache tzdata
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -o starter-kit ./cmd

FROM alpine:latest
RUN apk update && apk add --no-cache
WORKDIR /app
COPY --from=builder /app/ /app/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=Europe/Istanbul
ENTRYPOINT [ "/app/starter-kit"]