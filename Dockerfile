FROM golang:latest as build
RUN go get github.com/eclipse/paho.mqtt.golang && \
    go get github.com/mxk/go-imap/imap && \
    go get github.com/gorilla/mux
ADD . /go/src/github.com/kapitanov/miot-arc-lamp
WORKDIR /go/src/github.com/kapitanov/miot-arc-lamp
RUN go get
RUN CGO_ENABLED=0 go build -o miot-arc-lamp . 

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=build /go/src/github.com/kapitanov/miot-arc-lamp/miot-arc-lamp /app/miot-arc-lamp
COPY --from=build /go/src/github.com/kapitanov/miot-arc-lamp/www /app/www
EXPOSE 3000
WORKDIR /app
CMD ["/app/miot-arc-lamp"]