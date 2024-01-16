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
// versin 1.3.3
// if realm is emp, do nothing in function updateUser() in restapi.go
// version 1.5 [2022-06-21]
// auth server name changed : store.local -> 02.store.local, partner.local -> 03.partner.local

// version 1.6 [2024-01-10]
// EMP-GOTP realm add
// version 1.7 [2024-01-11]
// accountExpires add

// version 1.8 [2024-01-11]
// seperate config daysEMP for EMP-GOTP and days for others
// add disabled users for accountExpires

import (
	"context"
	"log"
	"os"
	"time"
)

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

	fpLog, err := os.OpenFile(logpath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer fpLog.Close()
	log.SetOutput(fpLog)

	scheduleTime := configuration.ScheduleTime

	//////////////////////
	// scheduledTask()
	// if scheduleTime <= 24 {
	// 	time.Sleep(time.Minute * 10)
	// 	os.Exit(0)
	// }
	//////////////////////

	for {
		now := time.Now()
		log.Println("firtst loop checked", now)
		if now.Hour() == scheduleTime-1 {
			for {
				now = time.Now()
				if now.Minute() >= 50 {
					for {
						now = time.Now()
						if now.Minute() == 0 {
							scheduledTask()
							goto SLEEP
						}
						time.Sleep(time.Minute * 1)
					}
				}
				time.Sleep(time.Minute * 10)
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
