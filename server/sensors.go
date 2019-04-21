package main

import (
	"fmt"
	"net/http"
	"strings"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/go-kivik/couchdb" // The CouchDB driver
	"github.com/gomodule/redigo/redis"
	"github.com/labstack/echo/v4"
)

//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

var sensorHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	st := strings.Split(msg.Topic(), "/")
	deviceType := st[1]
	deviceID := st[2]
	sensorType := st[3]

	_, err := DBGet(deviceType, deviceID)
	if err != nil {
		doc := make(map[string]interface{})
		doc[sensorType] = true
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

	n, err := r.Do("SET", fmt.Sprintf("%s:%s:%s ", deviceType, deviceID, sensorType), msg.Payload())
	if err != nil {
		fmt.Printf("Redis: %v\n", n)
	}
}

var opts *MQTT.ClientOptions
var c MQTT.Client

var sensorTopic string = "device/+/+/+"

//InitSensorService call in main
func InitSensorService(e *echo.Echo) {

	opts = MQTT.NewClientOptions().AddBroker("tcp://omaraa.ddns.net:1883")
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

	e.GET("/redis", redisTest)
}

var connectionLostHandler MQTT.ConnectionLostHandler = func(c MQTT.Client, err error) {
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

//DisconnectSensorService defer in main
func DisconnectSensorService() {
	//unsubscribe from /go-mqtt/sample
	if token := c.Unsubscribe(sensorTopic); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	c.Disconnect(250)
}

func redisTest(c echo.Context) error {
	r, err := redis.Dial("tcp", ":6379")
	if err != nil {
		panic(err)
	}
	defer r.Close()

	n, _ := r.Do("PING")

	return c.String(http.StatusOK, n.(string))
}

func PublishColor(deviceType string, deviceID string, color string) {
	c.Publish(fmt.Sprintf("device/%s/%s/color", deviceType, deviceID), 0, false, color)
}
