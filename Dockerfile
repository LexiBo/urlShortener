FROM golang:1.17.3-alpine3.14 as builder

ADD . /app
WORKDIR /app
RUN go mod tidy
RUN go build -o app cmd/main.go

FROM alpine:3.14
COPY --from=builder /app .
EXPOSE 8080
ENTRYPOINT ["./app"]