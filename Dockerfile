FROM golang:latest as build
RUN go get github.com/eclipse/paho.mqtt.golang && \
    go get github.com/mxk/go-imap/imap && \
    go get github.com/gorilla/mux
ADD . /go/src/github.com/kapitanov/miot-arc-lamp
WORKDIR /go/src/github.com/kapitanov/miot-arc-lamp
RUN go get
RUN go build -o miot-arc-lamp . 

FROM alpine:latest
COPY --from=build /go/src/github.com/kapitanov/miot-arc-lamp/miot-arc-lamp /app/miot-arc-lamp
COPY --from=build /go/src/github.com/kapitanov/miot-arc-lamp/www /app/www
CMD ["/app/miot-arc-lamp"]