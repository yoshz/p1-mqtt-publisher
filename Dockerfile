FROM golang AS BUILD
WORKDIR /src
ADD go.mod go.sum main.go /src/
RUN go get
RUN go build -ldflags "-linkmode external -extldflags -static" -o p1-mqtt-publisher

FROM scratch
COPY --from=build /src/p1-mqtt-publisher /
ENV SERIAL_DEVICE /dev/ttyUSB0
CMD ["/p1-mqtt-publisher"]
