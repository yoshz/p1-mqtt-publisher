# P1-MQTT-Publisher

A small Go application that reads the P1 serial device and publishes meter values as JSON over MQTT.

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
 * `MQTT_BROKER`: url of the MQTT broker (default is "tcp://localhost:1883")
 * `MQTT_USERNAME`: Username to authenticate to MQTT broker
 * `MQTT_PASSWORD`: Password to authenticate to MQTT broker
 * `MQTT_TOPIC`: MQTT Topic to publish to (default is "energy/meters")
