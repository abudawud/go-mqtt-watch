package main

type ConfigMQTT struct {
    Server  	 	string
	Topic			string
	ClientID		string
	QOS				int
}

type ConfigMail struct {
    Server  	 	string
	Port			int
	Username		string
	Password		string
}