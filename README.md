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

## JSON message

Example message:
```json
{
    "time": "2022-01-02T03:04:05.678910123+01:00",
    "location": "home",
    "powerDraw": 532,
    "powerMeter1": 123456,
    "powerMeter2": 78901,
    "gasMeter": 56723
}
```

Power is in `watt`.
Gas is in `dm3`.

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

## Home Assistant

To make the meter data available in Home Assistant, enable the MQTT integration
and add the following lines to the `configuration.yaml`:

```yaml
mqtt:
  sensor:
    - state_topic: "energy/meters"
      name: "Current Power Draw"
      object_id: power_draw
      expire_after: 60
      device_class: energy
      state_class: measurement
      unit_of_measurement: "w"
      value_template: "{{ value_json['powerDraw'] }}"

    - state_topic: "energy/meters"
      name: "Power Meter Tariff 1"
      object_id: power_meter_tariff_1
      expire_after: 60
      device_class: energy
      state_class: total_increasing
      unit_of_measurement: "kWh"
      value_template: "{{ value_json['powerMeter1'] / 1000 }}"

    - state_topic: "energy/meters"
      name: "Power Meter Tariff 2"
      object_id: power_meter_tariff_2
      expire_after: 60
      device_class: energy
      state_class: total_increasing
      unit_of_measurement: "kWh"
      value_template: "{{ value_json['powerMeter2'] / 1000 }}"

    - state_topic: "energy/meters"
      name: "Gas Meter"
      object_id: gas_meter
      expire_after: 60
      device_class: energy
      state_class: total_increasing
      unit_of_measurement: "m³"
      value_template: "{{ value_json['gasMeter'] / 1000 }}"
```
