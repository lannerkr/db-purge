package main

// version 1.1
// old user distinguished
// api put status to disable
// load configuration from conf.json
// mongodb 4.2 compatable
// after update pulse status, write status to mongodb
// logging add
// version 1.2
// when user doesn't exist in pulse, put user status to deleted
// realm store/partner applied
// version 1.3
// scheduling add
// version 1.3.2
// second loop condition fix : "now.Minute() = 50" --> "now.Minute() >= 50"

import (
	"context"
	"log"
	"os"
	"time"
)

// type Configuration struct {
// 	MongoDB     string
// 	CheckDays   int
// 	PulseUri    string
// 	PulseApiKey string
// }

// var configuration Configuration

type DisabledUsers struct {
	users  string
	status string
}

func init() {
	conf := os.Args[1]
	config(conf)
}
func main() {

	logpath := configuration.LogPath
	//logpath := "./dbpurge.log"
	fpLog, err := os.OpenFile(logpath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer fpLog.Close()
	log.SetOutput(fpLog)
	// log.Println(configuration)
	//Logger = log.New(fpLog, "", log.LstdFlags)

	// configuration = Configuration{
	// 	"mongoadmin:secret@192.168.0.157:27017",
	// 	90,
	// 	"https://192.168.0.174",
	// 	"ptN9OdCYWXnAwAgOmTdRPF9UfMXBqzQxkBy1OmpCwmw=",
	// }
	scheduleTime := configuration.ScheduleTime
	for {
		now := time.Now()
		log.Println("firtst loop checked", now)
		if now.Hour() == scheduleTime-1 {
			for {
				now = time.Now()
				// log.Println("second loop checked", now)
				if now.Minute() >= 50 {
					for {
						now = time.Now()
						// log.Println("third loop checked", now)
						if now.Minute() == 0 {
							scheduledTask()
							goto SLEEP
						}
						time.Sleep(time.Minute * 1)
						// fmt.Println("tick 1 min")
					}
				}
				time.Sleep(time.Minute * 10)
				// fmt.Println("tick 10 min")
			}
		}
	SLEEP:
		time.Sleep(time.Hour * 1)
	}
}

func scheduledTask() {
	log.Println("[schedule] scheduledTask DB purge Task is executed")

	client := connectdb()
	defer client.Disconnect(context.TODO())
	coll := getColl(client)

	purgedb(coll)
	oldusers(coll)
}
