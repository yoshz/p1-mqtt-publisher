package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/tarm/serial"
)

var (
	mqttBroker      string = getEnv("MQTT_BROKER", "tcp://localhost:1883")
	mqttUsername    string = getEnv("MQTT_USERNAME", "")
	mqttPassword    string = getEnv("MQTT_PASSWORD", "")
	mqttTopic       string = getEnv("MQTT_TOPIC", "energy/meters")
	serialDevice    string = os.Getenv("SERIAL_DEVICE")
	publishInterval int64  = 60
	reader          *bufio.Reader
	message         EnergyMeterMessage
)

type EnergyMeterMessage struct {
	PowerDraw   float64 `json:"powerDraw"`
	PowerMeter1 float64 `json:"powerMeter1"`
	PowerMeter2 float64 `json:"powerMeter2"`
	GasMeter    float64 `json:"gasMeter"`
}

func main() {

	if serialDevice != "" {
		fmt.Println("gonna use serial device")
		config := &serial.Config{Name: serialDevice, Baud: 115200}

		usb, err := serial.OpenPort(config)
		if err != nil {
			fmt.Printf("Could not open serial interface: %s", err)
			return
		}

		reader = bufio.NewReader(usb)
	} else {
		fmt.Println("Using example file")
		file, err := os.Open("examples/fulllist.txt")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		reader = bufio.NewReader(file)
	}

	go reading(reader)

	opts := MQTT.NewClientOptions()
	opts.SetClientID("p1-mqtt-publisher")
	opts.AddBroker(mqttBroker)
	if mqttUsername != "" {
		opts.SetUsername(mqttUsername)
	}
	if mqttPassword != "" {
		opts.SetPassword(mqttPassword)
	}

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Println("Connected to MQTT broker")

	// sleeping 10 seconds to prevent uninitialized scrapes
	time.Sleep(10 * time.Second)

	fmt.Println("now serving metrics")

	publish(client)
}

func publish(client MQTT.Client) {
	for {
		payload, err := json.Marshal(message)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(payload[:]))

		if token := client.Publish(mqttTopic, 0, false, payload); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}

		time.Sleep(time.Duration(publishInterval) * time.Second)
	}
}

func reading(source io.Reader) {
	var line string
	for {
		rawLine, err := reader.ReadBytes('\x0a')
		if err != nil {
			fmt.Println(err)
			return
		}
		line = string(rawLine[:])
		if strings.HasPrefix(line, "1-0:1.8.1") {
			tmpVal, err := strconv.ParseFloat(line[10:20], 64)
			if err != nil {
				fmt.Println(err)
				continue
			}
			message.PowerMeter1 = tmpVal
		} else if strings.HasPrefix(line, "1-0:1.8.2") {
			tmpVal, err := strconv.ParseFloat(line[10:20], 64)
			if err != nil {
				fmt.Println(err)
				continue
			}
			message.PowerMeter2 = tmpVal
		} else if strings.HasPrefix(line, "0-1:24.2.1") {
			tmpVal, err := strconv.ParseFloat(line[26:35], 64)
			if err != nil {
				fmt.Println(err)
				continue
			}
			message.GasMeter = tmpVal
		} else if strings.HasPrefix(line, "1-0:1.7.0") {
			tmpVal, err := strconv.ParseFloat(line[10:16], 64)
			if err != nil {
				fmt.Println(err)
				continue
			}
			message.PowerDraw = tmpVal
		}
		if serialDevice == "" {
			time.Sleep(200 * time.Millisecond)
		}
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
