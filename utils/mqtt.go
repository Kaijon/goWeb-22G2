package utils

import (
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	exClient mqtt.Client
	inClient mqtt.Client
	exChoke  chan [2]string
	inChoke  chan [2]string
}

const (
	MQTT_EXTERNAL_TOPIC_REC_CA42A = "status/recording/ca42a"
	MQTT_EXTERNAL_TOPIC_REC_WEB   = "status/recording/web"
	MQTT_EXTERNAL_TOPIC_REC       = "status/recording/#"
	MQTT_INTERNAL_TOPIC_EVENT     = "event"
	MQTT_INTERNAL_TOPIC_REC       = "status/recording/#"
	MQTT_INTERNAL_TOPIC_REC_CA42A = "status/recording/ca42a"
	MQTT_EXTERNAL_BROKER_URI      = "tcp://192.168.5.8:1883"
	MQTT_INTERNAL_BROKER_URI      = "tcp://127.0.0.1:8083"
	MQTT_EXTERNAL_CLIENT_ID       = "MqttEx"
	MQTT_INTERNAL_CLIENT_ID       = "MqttIn"
)

var MqttClient *Client
var exChoke = make(chan [2]string)
var inChoke = make(chan [2]string)

func MqttNewClient() *Client {
	c := &Client{}
	//c.exClient = createExClient(MQTT_EXTERNAL_BROKER_URI, MQTT_EXTERNAL_CLIENT_ID)
	c.inClient = createInClient(MQTT_INTERNAL_BROKER_URI, MQTT_INTERNAL_CLIENT_ID)
	return c
}

func (m *Client) Publish(clientId, topic, data string) {
	if clientId == MQTT_EXTERNAL_CLIENT_ID {
		log.Printf("%s->%s(VHA-10)", MQTT_EXTERNAL_CLIENT_ID, data)
		m.exClient.Publish(topic, 0, false, data)
	} else if clientId == MQTT_INTERNAL_CLIENT_ID {
		log.Printf("%s->%s(LOCAL)", MQTT_INTERNAL_CLIENT_ID, data)
		m.inClient.Publish(topic, 0, false, data)
	} else {
		log.Println("Not identify Mqtt client")
	}
}

func createExClient(brokerIp, id string) mqtt.Client {
	opts := mqtt.NewClientOptions().AddBroker(brokerIp).SetClientID(id)
	opts.SetDefaultPublishHandler(monitorEx)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Println(token.Error())
	}

	if token := client.Subscribe(MQTT_EXTERNAL_TOPIC_REC, 0, mqttExHandler); token.Wait() && token.Error() != nil {
		log.Println(token.Error())
	}

	return client
}

func createInClient(brokerIp, id string) mqtt.Client {
	opts := mqtt.NewClientOptions().AddBroker(brokerIp).SetClientID(id)
	opts.SetDefaultPublishHandler(monitorIn)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		(token.Error())
	}

	if token := client.Subscribe(MQTT_INTERNAL_TOPIC_REC, 0, mqttInHandler); token.Wait() && token.Error() != nil {
		log.Println(token.Error())
	}

	return client
}

func monitorEx(c mqtt.Client, msg mqtt.Message) {
	exChoke <- [2]string{msg.Topic(), string(msg.Payload())}
}

func monitorIn(c mqtt.Client, msg mqtt.Message) {
	inChoke <- [2]string{msg.Topic(), string(msg.Payload())}
}

func mqttInHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Recv topic: %s, data: %s", msg.Topic(), msg.Payload())
	//var data JsonEvent
	//err := json.Unmarshal(msg.Payload(), &data)
	//if err != nil {
	//	log.Println("JSON Unmarshall failed")
	//}
	//log.Println("\tSource: ", data.Source)
	//log.Println("\tParams: ", data.Params.Cmd)
	payload := string(msg.Payload())
	switch msg.Topic() {
	case MQTT_INTERNAL_TOPIC_REC_CA42A:
		if payload == "start" {
			log.Println("daemon rec start done")
			//MqttClient.Publish(MQTT_INTERNAL_CLIENT_ID, MQTT_EXTERNAL_TOPIC_REC_WEB, "{\"rec\":\"on_done\"}")
			//WsSendMessageToClients("on_done")
		} else if payload == "stop" {
			log.Println("daemon rec start done")
			//MqttClient.Publish(MQTT_INTERNAL_CLIENT_ID, MQTT_EXTERNAL_TOPIC_REC_WEB, "{\"rec\":\"off_done\"}")
			//WsSendMessageToClients("off_done")
		} else {
			log.Printf("Unknow cmd=%s", payload)
		}

	default:
		log.Println("Source Module not identify")
	}

	/*
		if err := c.ShouldBindJSON(&jsonData); err != nil {
			log.Println(err.Error())
			c.JSON(400, gin.H{"error": err.Error()})
			//return
		}
		log.Println("MQTT Packet: ")
		log.Println("\tTopic: ", jsonData.Topic)
		log.Println("\tSource: ", jsonData.Source)
		log.Println("\tDest: ", jsonData.Destination)
		log.Println("\tRoom: ", jsonData.Room)
		log.Println("\tParams: ", jsonData.Cmd)
	*/
}

func mqttExHandler(client mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()
	value := string(msg.Payload())

	switch topic {
	case MQTT_EXTERNAL_TOPIC_REC_CA42A:
		if value == "start" {
			log.Println("REC ON DONE")
			//MqttClient.Publish(MQTT_INTERNAL_CLIENT_ID, "REC_ON")
		} else if value == "stop" {
			log.Println("REC OFF DONE")
			//MqttClient.Publish(MQTT_INTERNAL_CLIENT_ID, "REC_OFF")
		} else {
			log.Printf("Unknow cmd=%s", value)
		}
	default:
		log.Printf("topic: %s, message: %s", msg.Topic(), msg.Payload())
		log.Println("Source Module not identify")
	}
}

func StartMqttInLoop() {
	log.Println("MqttIn Monitor start")
	for {
		incoming := <-inChoke
		log.Printf("Recv:Topic:%s, Msg:%s\n", incoming[0], incoming[1])
	}
}

func StartMqttExLoop() {
	log.Println("MqttEx Monitor start")
	for {
		incoming := <-exChoke
		log.Printf("Recv:Topic:%s, Msg:%s\n", incoming[0], incoming[1])
	}
}

func init() {
	MqttClient = MqttNewClient()
}
