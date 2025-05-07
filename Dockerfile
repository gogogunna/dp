FROM golang:1.24 AS my_app

WORKDIR /app

COPY . ./

RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o server

FROM alpine:latest

COPY --from=my_app /app/server ./

COPY --from=certs_context fullchain1.pem privkey1.pem /etc/ssl/certs/

EXPOSE 443

CMD ["./server"]
