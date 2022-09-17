package main

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
)

var (
	mqttBroker      string = getEnv("MQTT_BROKER", "tcp://localhost:1883")
	mqttUsername    string = getEnv("MQTT_USERNAME", "")
	mqttPassword    string = getEnv("MQTT_PASSWORD", "")
	mqttTopic       string = getEnv("MQTT_TOPIC", "energy/meters")
	serialDevice    string = os.Getenv("SERIAL_DEVICE")
	location        string = getEnv("LOCATION", "home")
	publishInterval int64  = 60
	reader          *bufio.Reader
	message         EnergyMeterMessage
)

type EnergyMeterMessage struct {
	Time        time.Time `json:"time"`
	Location    string    `json:"location"`
	PowerDraw   int64     `json:"powerDraw"`
	PowerMeter1 int64     `json:"powerMeter1"`
	PowerMeter2 int64     `json:"powerMeter2"`
	GasMeter    int64     `json:"gasMeter"`
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	message.Location = location

	if serialDevice != "" {
		log.Printf("Using serial device: %s", serialDevice)
		config := &serial.Config{Name: serialDevice, Baud: 115200}

		usb, err := serial.OpenPort(config)
		if err != nil {
			log.Fatalf("Could not open serial device: %s", err)
		}

		reader = bufio.NewReader(usb)
	} else {
		log.Println("Using example file")
		file, err := os.Open("examples/fulllist.txt")
		if err != nil {
			log.Fatalln(err)
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
		log.Fatalln(token.Error())
	}
	log.Println("Connected to MQTT broker")

	// sleeping 10 seconds to prevent uninitialized scrapes
	time.Sleep(10 * time.Second)

	log.Println("Now publishing metrics")

	publish(client)
}

func publish(client MQTT.Client) {
	for {
		message.Time = time.Now()
		payload, err := json.Marshal(message)
		if err != nil {
			log.Fatalln(err)
		}

		if token := client.Publish(mqttTopic, 1, false, payload); token.Wait() && token.Error() != nil {
			log.Fatalln(token.Error())
		}

		time.Sleep(time.Duration(publishInterval) * time.Second)
	}
}

func reading(source io.Reader) {
	var line string
	for {
		rawLine, err := reader.ReadBytes('\x0a')
		if err != nil {
			log.Errorln(err)
			return
		}
		line = string(rawLine[:])
		if strings.HasPrefix(line, "1-0:1.8.1") {
			tmpVal, err := strconv.ParseFloat(line[10:20], 64)
			if err != nil {
				log.Errorln(err)
				continue
			}
			message.PowerMeter1 = int64(tmpVal * 1000)
		} else if strings.HasPrefix(line, "1-0:1.8.2") {
			tmpVal, err := strconv.ParseFloat(line[10:20], 64)
			if err != nil {
				log.Errorln(err)
				continue
			}
			message.PowerMeter2 = int64(tmpVal * 1000)
		} else if strings.HasPrefix(line, "0-1:24.2.1") {
			tmpVal, err := strconv.ParseFloat(line[26:35], 64)
			if err != nil {
				log.Errorln(err)
				continue
			}
			message.GasMeter = int64(tmpVal * 1000)
		} else if strings.HasPrefix(line, "1-0:1.7.0") {
			tmpVal, err := strconv.ParseFloat(line[10:16], 64)
			if err != nil {
				log.Errorln(err)
				continue
			}
			message.PowerDraw = int64(tmpVal * 1000)
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
