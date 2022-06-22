package main

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	MongoDB      string
	CheckDays    int
	PulseUri     string
	PulseApiKey  string
	LogPath      string
	ScheduleTime int
}

var configuration Configuration

func config(conf string) {
	file, _ := os.Open(conf)
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration = Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Println("error:", err)
	}
	log.Println(configuration)
}
