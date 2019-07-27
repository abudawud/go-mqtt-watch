package main

type Sensors struct {
	ID        int   `json:"s"`
	Timestamp int64 `json:"t"`
	SensorA1  int32 `json:"ai1"`
	SensorAS1 int32 `json:"ai_st1"`
	SensorA2  int32 `json:"ai2"`
	SensorAS2 int32 `json:"ai_st2"`
}
