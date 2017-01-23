package handlers

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/alouca/gosnmp"
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

type DeviceInfo struct {
	DevIP string
	Value string
	Time  time.Time
}

type ConcretTask struct {
	DevIP  string
	DevOID string
	Repeat int64
}

func (self *ConcretTask) StartTask() {
	for {
		dev := SNMPCheckStart(self.DevIP, self.DevOID)
		fmt.Printf("|V:%s | D:%s | T:%s|\n", dev.Value, dev.DevIP, dev.Time)
		time.Sleep(time.Second * time.Duration(self.Repeat))
	}
}

func (self *TaskListCreator) CreatTask(devip string, devOID string, repeat int64) Task {
	var task Task
	task = &ConcretTask{devip, devOID, repeat}
	self.registerTask(task)
	return task
}

func SNMPCheckStart(deviceIP string, oid string) DeviceInfo {
	var dev DeviceInfo
	s, err := gosnmp.NewGoSNMP(deviceIP, "public", gosnmp.Version2c, 5)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := s.Get(oid)
	if err == nil {
		for _, val := range resp.Variables {
			fmt.Println(val)
			switch val.Type {
			case gosnmp.OctetString:
				dev.Value = val.Value.(string)
				dev.Time = time.Now()
				dev.DevIP = deviceIP
			case gosnmp.BitString:
				dev.Value = strconv.FormatUint(val.Value.(uint64), 10)
				dev.Time = time.Now()
				dev.DevIP = deviceIP
			case gosnmp.Counter64:
				dev.Value = strconv.FormatUint(val.Value.(uint64), 10)
				dev.Time = time.Now()
				dev.DevIP = deviceIP
			}
		}
	}
	return dev
}
