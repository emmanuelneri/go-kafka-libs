FROM golang:1.14-alpine as builder
RUN apk update && apk add --no-cache git
WORKDIR build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/app ./cmd

FROM scratch
COPY --from=builder /bin/app /bin/app
ENV KAFKA_BOOTSTRAP_SERVERS "kafka:29092"
ENTRYPOINT ["/bin/app"]
EXPOSE 8082