package datasource

import (
	"../devices"
	"../handlers"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	config handlers.Config
)

func (self *MonitoringBase) CheckError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func (self *MonitoringBase) connectSession(user string, pass string, host string, port string) (*mgo.Session, error) {
	if len(port) == 0 {
		port = "27017"
	}
	url := "mongodb://" + user + ":" + pass + "@" + host + ":" + port
	return mgo.Dial(url)
}

func (self *MonitoringBase) sessionStart() (*mgo.Session, error) {
	config = handlers.GetConfigData() // load config
	return self.connectSession(config.Login, config.Password, config.Ip, config.Port)
}

func (self *MonitoringBase) loadData(table string, data *[]interface{}) {
	session, err := self.sessionStart()
	self.CheckError(err)
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(config.Base).C(table)
	var result []interface{}
	err = c.Find(bson.M{}).All(&result)
	self.CheckError(err)
	*data = result
}
//devicegroup
func (self *MonitoringBase) LoadDeviceGroup() []devices.DeviceGroup {
	var devgroupI []interface{}
	self.loadData("devicegroup", &devgroupI)
	var snmptemplate []devices.DeviceGroup

	for _, v := range devgroupI {
		var st devices.DeviceGroup
		bsonBytes, _ := bson.Marshal(v)
		bson.Unmarshal(bsonBytes, &st)
		snmptemplate = append(snmptemplate, st)
	}
	return snmptemplate
}
//netdev
func (self *MonitoringBase) LoadNetDevice() []devices.NetDev {
	var devgroupI []interface{}
	self.loadData("netdevice", &devgroupI)
	var netDev []devices.NetDev

	for _, v := range devgroupI {
		var st devices.NetDev
		bsonBytes, _ := bson.Marshal(v)
		bson.Unmarshal(bsonBytes, &st)
		netDev = append(netDev, st)
	}
	return netDev
}

func (self *MonitoringBase) insertData(table string, data interface{}) {
	session, err := self.sessionStart()
	self.CheckError(err)
	c := session.DB(config.Base).C(table)
	err = c.Insert(&data)
	self.CheckError(err)
}

func (self *MonitoringBase) WriteDeviceGroup(deviceName string) {
	dev_group := devices.DeviceGroup{bson.NewObjectId(), deviceName}
	self.insertData("devicegroup", dev_group)
}

func (self *MonitoringBase) WriteNetDev(name string, locate string, ip string, active bool, groupid bson.ObjectId) {
	net_dev := devices.NetDev{bson.NewObjectId(), name, locate, ip, active, groupid}
	self.insertData("netdevice", net_dev)
}

func (self *MonitoringBase) WriteOidList(name string, oid string, groupid int64, repeat int64) {
	oid_list := devices.OidList{bson.NewObjectId(), name, oid, groupid, repeat}
	self.insertData("oidlist", oid_list)
}
