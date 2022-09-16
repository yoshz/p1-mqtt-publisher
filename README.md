# P1-MQTT-Publisher

A small Go application that reads the P1 serial device connected to a SmartMeter and publishes meter values as JSON over MQTT.

## Local usage

Start MQTT broker:

```bash
docker-compose up -d
```

Install Go modules:
```bash
go get
```

Run the P1-MQTT-Publisher
```bash
go run main.go
```

## Environment variables

 * `SERIAL_DEVICE`: the device that needs to be read (default is /dev/ttyUSB0)
 * `LOCATION`: Location of the SmartMeter (default is "home")
 * `MQTT_BROKER`: url of the MQTT broker (default is "tcp://localhost:1883")
 * `MQTT_USERNAME`: Username to authenticate to MQTT broker
 * `MQTT_PASSWORD`: Password to authenticate to MQTT broker
 * `MQTT_TOPIC`: MQTT Topic to publish to (default is "energy/meters")

## Docker image

See https://hub.docker.com/r/yoshz/p1-mqtt-publisher/tags

## Kubernetes Deployment

```bash
# create namespace
kubectl create namespace p1

# install mosquitto
kubectl apply -n p1 -k deploy/mosquitto

# install p1-mqtt-publisher
kubectl apply -n p1 -k deploy/p1-mqtt-publisher
```
