package handlers

import (
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

	var res1 *gosnmp.SnmpPacket

	var res2 *gosnmp.SnmpPacket

	res1, _ = s.Get(oid)
	time.Sleep(time.Second * time.Duration(1))
	res2, _ = s.Get(oid)

	q1, _ := strconv.ParseInt(snmpDataParser(res1), 10, 64)
	q2, _ := strconv.ParseInt(snmpDataParser(res2), 10, 64)

	dev.Value = strconv.FormatInt((q2 - q1), 10)

	dev.DeviceID = deviceID
	dev.Time = time.Now()
	dev.DevIP = deviceIP
	dev.CheckType = checkType
	return dev
}

func snmpDataParser(data interface{}) string {
	var result string
	if r, ok := data.(*gosnmp.SnmpPacket); ok {
		for _, val := range r.Variables {
			switch val.Type {
			case gosnmp.OctetString:
				result = val.Value.(string)

			case gosnmp.BitString:
				result = strconv.FormatUint(val.Value.(uint64), 10)

			case gosnmp.Counter64:
				result = strconv.FormatUint(val.Value.(uint64), 10)

			}
		}
	}
	return result
}
