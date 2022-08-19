package models

type LogMessage struct {
	Level     string `bson:"level" json:"level"`
	Timestamp string `bson:"timestamp" json:"ts"`
	Caller    string `bson:"caller" json:"caller"`
	Message   string `bson:"error_message" json:"msg"`
}
