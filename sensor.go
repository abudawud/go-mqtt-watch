package main

type Sensors struct {
	Timestamp string `json:"t"`
	Sensor1   int    `json:"ai1,string"`
	Sensor2   int    `json:"ai2,string"`
	Sensor3   int    `json:"ai3,string"`
	Sensor4   int    `json:"ai4,string"`
}
