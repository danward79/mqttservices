package main

import (
	"flag"
	"fmt"

	"github.com/danward79/mqttservices"
)

func main() {
	//Command line variables
	topic := flag.String("t", "home/#", "Enter the topic to subscribe to.")
	mqttServer := flag.String("s", ":1883", "Enter the IP and Port of the MQTT Broker. e.g. 127.0.0.1:1883")
	flag.Parse()

	mqttClient := mqttservices.NewClient(*mqttServer)
	chSub := mqttClient.Subscribe(*topic)

	for {
		select {
		case m := <-chSub:
			fmt.Printf("%s\t\t%s\n", m.TopicName, m.Payload)

		}
	}
}
