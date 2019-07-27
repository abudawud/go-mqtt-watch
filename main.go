package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"encoding/json"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	gonfig "github.com/tkanos/gonfig"
)

var cfgMail ConfigMail
var MAIL_TO = []string{"warishafidz@gmail.com", "rochieirawan2405@gmail.com", "muhammadfadlullah9@gmail.com"}

const MAIL_FROM = "abudawud@sabinsolusi.com"
const MAIL_SUBJECT = "Gas sensor Notification!"

func onMessageReceived(client MQTT.Client, message MQTT.Message) {
	var sensors Sensors
	json.Unmarshal([]byte(message.Payload()), &sensors)

	if sensors.SensorA1 > 50 {
		body := fmt.Sprintf("Device dengan ID: <b>MAC_ADDR</b> kelebihan tekanan <br>")
		body = fmt.Sprintf("%s Sensor ID: <b>%d</b> <br>", body, sensors.ID)
		body = fmt.Sprintf("%s Tekanan saat ini: <b>%d</b> <br>", body, sensors.SensorA1)
		body = fmt.Sprintf("%s Lokasi Device: <b>LOKASI</b>", body)
		mail := Mail{
			MAIL_FROM,
			MAIL_TO,
			MAIL_SUBJECT,
			body,
		}
		PostMail(sensors.ID, mail)
	}
}

func main() {
	cfgMQTT := ConfigMQTT{}
	cfgMail = ConfigMail{}

	err := gonfig.GetConf("conf/mqtt.json", &cfgMQTT)
	err1 := gonfig.GetConf("conf/mail.json", &cfgMail)

	if err != nil && err1 != nil {
		panic(err)
	}

	fmt.Printf("Broker: %s\n", cfgMQTT.Server)
	fmt.Printf("Topic: %s\n", cfgMQTT.Topic)
	fmt.Printf("clientID: %s\n", cfgMQTT.ClientID)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	InitMailer(cfgMail)
	go MailSender()

	connOpts := MQTT.NewClientOptions().AddBroker(cfgMQTT.Server).SetClientID(cfgMQTT.ClientID).SetCleanSession(true)
	connOpts.OnConnect = func(c MQTT.Client) {
		if token := c.Subscribe(cfgMQTT.Topic, byte(cfgMQTT.QOS), onMessageReceived); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}

	client := MQTT.NewClient(connOpts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Printf("Connected to %s\n", cfgMQTT.Server)
	}

	<-c
}
