package main

import (
	"./webservice"

	"fmt"
)

func main() {

	go webservice.WebServer() // run web service




	/*base := datasource.MonitoringBase{}
	dev := base.LoadDeviceGroup()
	fmt.Println(dev)
*/
	/*	factory := new(handlers.TaskListCreator)
		tasks := []handlers.Task{
			factory.CreatTask("192.168.88.1", ".1.3.6.1.2.1.1.1.0", 5),
			factory.CreatTask("192.168.88.1", ".1.3.6.1.2.1.31.1.1.1.7.2", 5),
			factory.CreatTask("192.168.88.1", ".1.3.6.1.2.1.31.1.1.1.11.2", 5),
		}
		for _, t := range tasks {
			go t.StartTask()
		}*/
	var intput string
	fmt.Scanln(&intput)
}
