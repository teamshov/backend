package main

import (
	"fmt"
	"strings"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/go-kivik/couchdb" // The CouchDB driver
	"github.com/gomodule/redigo/redis"
	"github.com/tarm/serial"
)

//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

var s *serial.Port

var sensorHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	st := strings.Split(msg.Topic(), "/")
	deviceType := st[1]
	deviceID := st[2]
	sensorType := st[3]

	if sensorType == "color" {
		return
	}

	_, err := DBGet(deviceType, deviceID)
	if err != nil {
		doc := make(map[string]interface{})
		doc[sensorType] = true
		doc["xpos"] = 0
		doc["ypos"] = 0

		if devicesSettings[deviceType].(map[string]interface{})["color"].(bool) {
			doc["color"] = true
		}

		err := DBPut(deviceType, deviceID, doc)
		if err != nil {
			fmt.Println("Failed to put device")
		}
	}

	r, err := redis.Dial("tcp", ":6379")
	if err != nil {
		panic(err)
	}
	defer r.Close()

	path := fmt.Sprintf("device:%s:%s:dangerlevel", deviceType, deviceID)
	n, err := r.Do("SET", path, msg.Payload())
	if err != nil {
		fmt.Printf("Redis SET: %v\n", n)
	}

}

var opts *MQTT.ClientOptions
var c MQTT.Client
var devicesSettings map[string]interface{}

var sensorTopic string = "device/+/+/+"

//InitSensorService call in main
func InitSensorService() {

	devicesSettings, _ = DBGet("system", "devices")

	opts = MQTT.NewClientOptions().AddBroker("tcp://localhost:1883")
	opts.SetClientID("backend")
	opts.SetDefaultPublishHandler(f)
	opts.SetConnectionLostHandler(connectionLostHandler)

	//create and start a client using the above ClientOptions
	c = MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := c.Subscribe(sensorTopic, 0, sensorHandler); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

}

var connectionLostHandler MQTT.ConnectionLostHandler = func(c MQTT.Client, err error) {
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

//DisconnectSensorService defer in main
func DisconnectSensorService() {
	if token := c.Unsubscribe(sensorTopic); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	c.Disconnect(250)
}

func PublishColor(deviceType string, deviceID string, color string) {
	c.Publish(fmt.Sprintf("device/%s/%s/color", deviceType, deviceID), 0, false, color)
}
