package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"encoding/json"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	gonfig "github.com/tkanos/gonfig"
)

var cfgMail ConfigMail
var MAIL_TO = []string{"warishafidz@gmail.com", "rochieirawan2405@gmail.com", "suwondo@smartlintas.com", "muhammadfadlullah9@gmail.com"}

const MAIL_FROM = "abudawud@sabinsolusi.com"
const MAIL_SUBJECT = "Gas sensor Notification!"

func CreateMail(sValue int, sId string, devId string) (mail Mail) {
	body := fmt.Sprintf("Device dengan ID: <b>%s</b> kelebihan tekanan <br>", devId)
	body = fmt.Sprintf("%s Sensor ID: <b>%s</b> <br>", body, sId)
	body = fmt.Sprintf("%s Tekanan saat ini: <b>%d</b> <br>", body, sValue)
	body = fmt.Sprintf("%s Lokasi Device: <b>LOKASI</b>", body)
	mail = Mail{
		MAIL_FROM,
		MAIL_TO,
		MAIL_SUBJECT,
		body,
	}

	return mail
}

func onMessageReceived(client MQTT.Client, message MQTT.Message) {
	var sensors Sensors
	id := strings.Replace(message.Topic(), "/data", "", 1)
	json.Unmarshal([]byte(message.Payload()), &sensors)

	fmt.Printf("Topic: %s, msg: %s\n", id, message.Payload())
	fmt.Println(sensors.Sensor1)
	if sensors.Sensor1 > 50 {
		mail := CreateMail(sensors.Sensor1, "1", id)

		PostMail(1, mail)
	}
	if sensors.Sensor2 > 50 {
		mail := CreateMail(sensors.Sensor2, "2", id)

		PostMail(2, mail)
	}
	if sensors.Sensor3 > 50 {
		mail := CreateMail(sensors.Sensor3, "3", id)

		PostMail(3, mail)
	}
	if sensors.Sensor4 > 50 {
		mail := CreateMail(sensors.Sensor4, "4", id)

		PostMail(3, mail)
	}
}

func main() {
	cfgMQTT := ConfigMQTT{}
	cfgMail = ConfigMail{}

	err := gonfig.GetConf("conf.d/mqtt.json", &cfgMQTT)
	err1 := gonfig.GetConf("conf.d/mail.json", &cfgMail)

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
