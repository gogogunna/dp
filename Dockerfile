FROM golang:1.24 AS my_app

WORKDIR /app

COPY . ./

RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o server

FROM alpine:latest

COPY --from=my_app /app/server ./

EXPOSE 80

CMD ["./server"]
