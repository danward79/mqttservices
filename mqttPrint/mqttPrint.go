package main

import (
	"flag"
	"log"

	"github.com/danward79/mqttservices"
	proto "github.com/huin/mqtt"
)

func main() {
	//Command line variables
	topic := flag.String("t", "home/#", "Enter the topic to subscribe to.")
	mqttServer := flag.String("s", ":1883", "Enter the IP and Port of the MQTT Broker. e.g. 127.0.0.1:1883")
	flag.Parse()

	mqttClient := mqttservices.NewClient(*mqttServer)

	chSub := mqttClient.Subscribe([]proto.TopicQos{{
		Topic: *topic,
		Qos:   proto.QosAtMostOnce,
	}})

	for m := range chSub {
		log.Printf("%s\t\t%s\n", m.TopicName, m.Payload)
	}
}
