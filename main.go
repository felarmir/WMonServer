package main

import (
	"fmt"

	"./datasource"
	"./handlers"
	"./webservice"
)

func main() {
	go webservice.WebServer() // run web service

	base := datasource.MonitoringBase{}
	devices := base.LoadNetDevice()

	oids := base.LoadOIDListBy(devices[0].Groupid)

	factory := new(handlers.TaskListCreator)
	tasks := []handlers.Task{}
	for _, dev := range devices {
		for _, oid := range oids {
			tasks = append(tasks, factory.CreatTask(dev.ID, dev.IP, oid.Oid, oid.Name, oid.Repeat))
		}
	}

	for _, t := range tasks {
		go t.StartTask()
	}

	var intput string
	fmt.Scanln(&intput)
}
