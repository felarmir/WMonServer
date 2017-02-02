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

// munction for check error
func (self *MonitoringBase) CheckError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// function for start session with mongodb server
func (self *MonitoringBase) sessionStart() (*mgo.Session, error) {
	config = handlers.GetConfigData() // load config

	if len(config.Port) == 0 {
		config.Port = "27017"
	}
	url := "mongodb://" + config.Login + ":" + config.Password + "@" + config.Ip + ":" + config.Port
	return mgo.Dial(url)
}

// Load data by table name and return intrface
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

//Load Datat from devicegroup and cacting interface to struct DeviceGroup array
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

//Load Data from netdevice table and casting interface to NetDev struct
func (self *MonitoringBase) LoadNetDevice() []devices.NetDevice {
	var devgroupI []interface{}
	self.loadData("netdevice", &devgroupI)
	var netDev []devices.NetDevice

	for _, v := range devgroupI {
		var st devices.NetDevice
		bsonBytes, _ := bson.Marshal(v)
		bson.Unmarshal(bsonBytes, &st)
		netDev = append(netDev, st)
	}
	return netDev
}

// Insert interface Data to table
func (self *MonitoringBase) insertData(table string, data interface{}) {
	session, err := self.sessionStart()
	self.CheckError(err)
	c := session.DB(config.Base).C(table)
	err = c.Insert(&data)
	self.CheckError(err)
}

// Delete Data Row in table by rowID
func (self *MonitoringBase) DeleteDataRow(table string, rowID string) {
	session, err := self.sessionStart()
	self.CheckError(err)
	c := session.DB(config.Base).C(table)
	err = c.Remove(bson.M{"_id": bson.ObjectIdHex(rowID)})
	self.CheckError(err)
}

//Update Data Row in table by ID
func (self *MonitoringBase) UpdateDataRow(table string, rowID string, newData bson.M) {
	session, err := self.sessionStart()
	self.CheckError(err)
	c := session.DB(config.Base).C(table)
	rowIdent := bson.M{"_id": bson.ObjectIdHex(rowID)}
	err = c.Update(rowIdent, newData)
	self.CheckError(err)
}

// Page Data loader and casting to Page Array
func (self *MonitoringBase) LoadMonitoringPages() []MonitoringPages {
	var result []interface{}
	self.loadData("pages", &result)
	var pages []MonitoringPages
	for _, v := range result {
		var p MonitoringPages
		bsonBytes, _ := bson.Marshal(v)
		bson.Unmarshal(bsonBytes, &p)
		pages = append(pages, p)
	}
	return pages
}

// Load Single Page Data by Page ID
func (self *MonitoringBase) LoadMonitoringPage(pageID string) MonitoringPages {
	session, err := self.sessionStart()
	self.CheckError(err)
	c := session.DB(config.Base).C("pages")
	var result interface{}
	err = c.Find(bson.M{"_id": bson.ObjectIdHex(pageID)}).One(&result)
	var page MonitoringPages

	bsonBytes, _ := bson.Marshal(result)
	bson.Unmarshal(bsonBytes, &page)
	return page
}

//Load menu. Load data and cast interface to MenuGroups
func (self *MonitoringBase) MenuGroupsList() []MenuGroups {
	var result []interface{}
	self.loadData("menugroup", &result)
	var menu []MenuGroups

	for _, v := range result{
		var tmp MenuGroups
		bsonBytes, _ := bson.Marshal(v)
		bson.Unmarshal(bsonBytes, &tmp)
		menu = append(menu, tmp)
	}
	return menu
}

func (self *MonitoringBase) ChildMenuList() []ChildMenu {
	var result []interface{}
	self.loadData("childmenu", &result)
	var child []ChildMenu
	for _, v := range result {
		var tmp ChildMenu
		bsBytes, _ := bson.Marshal(v)
		bson.Unmarshal(bsBytes, &tmp)
		child = append(child, tmp)
	}
	return child
}


// Write Device Group list
func (self *MonitoringBase) WriteDeviceGroup(deviceName string) {
	dev_group := devices.DeviceGroup{bson.NewObjectId(), deviceName}
	self.insertData("devicegroup", dev_group)
}

// Write Network device list
func (self *MonitoringBase) WriteNetDev(name string, locate string, ip string, active bool, groupid bson.ObjectId) {
	net_dev := devices.NetDevice{bson.NewObjectId(), name, locate, ip, active, groupid}
	self.insertData("netdevice", net_dev)
}

// Write OID List
func (self *MonitoringBase) WriteOidList(name string, oid string, groupid int64, repeat int64) {
	oid_list := devices.OidList{bson.NewObjectId(), name, oid, groupid, repeat}
	self.insertData("oidlist", oid_list)
}

// Write Menu Group
func (self *MonitoringBase) WriteMenuGroupList(menuTitle string, pageid string, table string) {
	menuGroup := MenuGroups{bson.NewObjectId(), menuTitle, pageid}
	self.insertData(table, menuGroup)
}

func (self *MonitoringBase) WriteMonitoringPage(pageName string, pageWg string, pageTable string, table string) {
	page := MonitoringPages{bson.NewObjectId(), pageName, pageWg, pageTable}
	self.insertData(table, page)
}

func (self *MonitoringBase) WriteChildMenu(title string, parent string, pageid string, table string) {
	child := ChildMenu{bson.NewObjectId(), title, parent, pageid}
	self.insertData(table, child)
}
