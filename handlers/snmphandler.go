package handlers

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"../datasource"
	"../devices"
	"github.com/alouca/gosnmp"
	"gopkg.in/mgo.v2/bson"
)

type Creator interface {
	CreatTask(devip string, devOID string, repeat int64) Task
	registerTask(task Task)
}

type Task interface {
	StartTask()
}

type TaskListCreator struct {
	taskList []*Task
}

func (self *TaskListCreator) registerTask(task Task) {
	self.taskList = append(self.taskList, &task)
}

type ConcretTask struct {
	DeviceID  bson.ObjectId
	DevIP     string
	DevOID    string
	CheckType string
	Repeat    int64
}

func (self *ConcretTask) StartTask() {
	for {
		dev := SNMPCheckStart(self.DeviceID, self.DevIP, self.DevOID, self.CheckType)
		mb := datasource.MonitoringBase{}
		mb.WriteDeviceStatus(dev.DeviceID, dev.DevIP, dev.Value, dev.CheckType, dev.Time)
		//fmt.Printf("|Value:%s | Device:%s | Time:%s|\n | Chek: %s", dev.Value, dev.DevIP, dev.Time, dev.CheckType)
		time.Sleep(time.Second * time.Duration(self.Repeat))
	}
}

func (self *TaskListCreator) CreatTask(deviceID bson.ObjectId, devip string, devOID string, checkType string, repeat int64) Task {
	var task Task
	task = &ConcretTask{deviceID, devip, devOID, checkType, repeat}
	self.registerTask(task)
	return task
}

func SNMPCheckStart(deviceID bson.ObjectId, deviceIP string, oid string, checkType string) devices.DeviceInfo {
	var dev devices.DeviceInfo
	s, err := gosnmp.NewGoSNMP(deviceIP, "public", gosnmp.Version2c, 5)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := s.Get(oid)
	if err == nil {
		for _, val := range resp.Variables {
			switch val.Type {
			case gosnmp.OctetString:
				dev.Value = val.Value.(string)

			case gosnmp.BitString:
				dev.Value = strconv.FormatUint(val.Value.(uint64), 10)

			case gosnmp.Counter64:
				dev.Value = strconv.FormatUint(val.Value.(uint64), 10)

			}
		}
	} else {
		fmt.Print(err)
	}
	dev.DeviceID = deviceID
	dev.Time = time.Now()
	dev.DevIP = deviceIP
	dev.CheckType = checkType
	return dev
}
