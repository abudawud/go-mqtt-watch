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
 
 func onMessageReceived(client MQTT.Client, message MQTT.Message) {
	var sensors Sensors
	json.Unmarshal([]byte(message.Payload()), &sensors)
	 fmt.Printf("Received message on topic: %s\nMessage: %s\n", message.Topic(), message.Payload())
	 fmt.Printf("Message: %d", sensors.S1)
 }
 
 func main() {
	 config := ConfigMQTT{}
	 err := gonfig.GetConf("config/mqtt.json", &config)

	 if err != nil {  panic(err) }

	 fmt.Printf("Broker: %s\n", config.Server)
	 fmt.Printf("Topic: %s\n", config.Topic)
	 fmt.Printf("clientID: %s\n", config.ClientID)

	 c := make(chan os.Signal, 1)
	 signal.Notify(c, os.Interrupt, syscall.SIGTERM)
 
	 connOpts := MQTT.NewClientOptions().AddBroker(config.Server).SetClientID(config.ClientID).SetCleanSession(true)
	 connOpts.OnConnect = func(c MQTT.Client) {
		 if token := c.Subscribe(config.Topic, byte(config.QOS), onMessageReceived); token.Wait() && token.Error() != nil {
			 panic(token.Error())
		 }
	 }
 
	 client := MQTT.NewClient(connOpts)
	 if token := client.Connect(); token.Wait() && token.Error() != nil {
		 panic(token.Error())
	 } else {
		 fmt.Printf("Connected to %s\n", config.Server)
	 }
 
	 <-c
 }