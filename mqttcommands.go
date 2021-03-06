// Package mqttservices provides an MQTT broker, topic subscription and publishing methods
package mqttservices

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	proto "github.com/huin/mqtt"
	"github.com/jeffallen/mqtt"
)

//Store client connection details
var ccPub *mqtt.ClientConn
var ccSub *mqtt.ClientConn

//MqttClient Subscripton struct
type MqttClient struct {
	Port string
}

//String returns details of the MqttClient
func (c *MqttClient) String() string {
	return fmt.Sprintf("MqttClient: IP %s", c.Port)
}

//NewClient declares a new broker
func NewClient(port string) *MqttClient {
	m := &MqttClient{Port: port}
	log.Println(m)
	return m
}

//Subscribe to MQTT Topic, takes topic as a string
func (c *MqttClient) Subscribe(tl []proto.TopicQos) chan *proto.Publish {

	if ccSub == nil {

		con, err := net.Dial("tcp", c.Port)
		gotError(err)

		ccSub = mqtt.NewClientConn(con)

		err = ccSub.Connect("", "")
		gotError(err)
	}

	ccSub.Subscribe(tl)

	return ccSub.Incoming
}

//Publish MQTT message, takes topic as a string, data as a byte array and retain flag as bool
func (c *MqttClient) Publish(topic string, data string, retain bool) {

	if ccPub == nil {
		con, err := net.Dial("tcp", c.Port)
		gotError(err)

		ccPub = mqtt.NewClientConn(con)

		err = ccPub.Connect("", "")
		gotError(err)
	}

	ccPub.Publish(&proto.Publish{
		Header: proto.Header{
			Retain: retain,
		},
		TopicName: topic,
		Payload:   proto.BytesPayload(data),
	})

}

//PublishMap a Map data entry takes a channel
func (c *MqttClient) PublishMap(chIn chan map[string]interface{}) {

	for ch := range chIn {
		for k := range ch {
			topic := generateTopic(ch["location"].(string))
			data := ""

			if k != "nodeid" {
				if k != "location" {
					switch value := ch[k].(type) {
					case string:
						data = value
					case int:
						data = strconv.Itoa(value)
					case bool:
						data = btos(value)
					case float64:
						data = strconv.FormatFloat(value, 'f', 2, 64)
					}

					c.Publish(topic+k, data, false)
				}
			}
		}
	}
}

//Generate location string
func generateTopic(s string) string {
	l := strings.Split(strings.ToLower(s), " ")

	loc := "home/"

	for _, v := range l {
		loc = loc + v + "/"
	}

	return loc
}

//Quick func to turn a bool to a string representation
func btos(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
