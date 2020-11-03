FROM golang:1.14 as app-builder
WORKDIR build
ADD . .
RUN go build -o /bin/app ./cmd

FROM golang:1.14
ENV KAFKA_BOOTSTRAP_SERVERS "kafka:29092"
COPY --from=app-builder /bin/app /bin/app
ENTRYPOINT ["/bin/app"]
EXPOSE 8080