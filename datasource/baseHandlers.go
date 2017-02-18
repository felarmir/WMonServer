package datasource

import (
	"fmt"

	"../devices"
	"../handlers"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	config handlers.Config
)

// munction for check error
func (mb *MonitoringBase) checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// function for start session with mongodb server
func (mb *MonitoringBase) sessionStart() (*mgo.Session, error) {
	config = handlers.GetConfigData() // load config

	if len(config.Port) == 0 {
		config.Port = "27017"
	}
	var url string
	if len(config.Login) == 0 {
		url = config.Ip
	} else {
		url = "mongodb://" + config.Login + ":" + config.Password + "@" + config.Ip + ":" + config.Port
	}
	return mgo.Dial(url)
}

// Load data by table name and return intrface
func (mb *MonitoringBase) loadData(table string, data *[]interface{}) {
	session, err := mb.sessionStart()
	mb.checkError(err)
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(config.Base).C(table)
	var result []interface{}
	err = c.Find(bson.M{}).All(&result)
	mb.checkError(err)
	*data = result
}

// Insert interface Data to table
func (mb *MonitoringBase) insertData(table string, data interface{}) {
	session, err := mb.sessionStart()
	mb.checkError(err)
	c := session.DB(config.Base).C(table)
	err = c.Insert(&data)
	mb.checkError(err)
}

//Load Datat from devicegroup and cacting interface to struct DeviceGroup array
func (mb *MonitoringBase) LoadDeviceGroup() []devices.DeviceGroup {
	var devgroupI []interface{}
	mb.loadData(DeviceGroupDBTable, &devgroupI)
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
func (mb *MonitoringBase) LoadNetDevice() []devices.NetDevice {
	var devgroupI []interface{}
	mb.loadData(NetDeviceDBTable, &devgroupI)
	var netDev []devices.NetDevice

	for _, v := range devgroupI {
		var st devices.NetDevice
		bsonBytes, _ := bson.Marshal(v)
		bson.Unmarshal(bsonBytes, &st)
		netDev = append(netDev, st)
	}
	return netDev
}

// Delete Data Row in table by rowID
func (mb *MonitoringBase) DeleteDataRow(table string, rowID string) {
	session, err := mb.sessionStart()
	mb.checkError(err)
	c := session.DB(config.Base).C(table)
	err = c.Remove(bson.M{"_id": bson.ObjectIdHex(rowID)})
	mb.checkError(err)
}

//Update Data Row in table by ID
func (mb *MonitoringBase) UpdateDataRow(table string, rowID string, newData bson.M) {
	session, err := mb.sessionStart()
	mb.checkError(err)
	c := session.DB(config.Base).C(table)
	rowIdent := bson.M{"_id": bson.ObjectIdHex(rowID)}
	err = c.Update(rowIdent, newData)
	mb.checkError(err)
}

// Page Data loader and casting to Page Array
func (mb *MonitoringBase) LoadMonitoringPages() []MonitoringPages {
	var result []interface{}
	mb.loadData(MonitorinPagesDBTable, &result)
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
func (mb *MonitoringBase) LoadMonitoringPage(pageID string) MonitoringPages {
	session, err := mb.sessionStart()
	mb.checkError(err)
	c := session.DB(config.Base).C(MonitorinPagesDBTable)
	var result interface{}
	err = c.Find(bson.M{"_id": bson.ObjectIdHex(pageID)}).One(&result)
	var page MonitoringPages

	bsonBytes, _ := bson.Marshal(result)
	bson.Unmarshal(bsonBytes, &page)
	return page
}

//Load menu. Load data and cast interface to MenuGroups
func (mb *MonitoringBase) MenuGroupsList() []MenuGroups {
	var result []interface{}
	mb.loadData(MenuGroupDBTable, &result)
	var menu []MenuGroups

	for _, v := range result {
		var tmp MenuGroups
		bsonBytes, _ := bson.Marshal(v)
		bson.Unmarshal(bsonBytes, &tmp)
		menu = append(menu, tmp)
	}
	return menu
}

func (mb *MonitoringBase) ChildMenuList() []ChildMenu {
	var result []interface{}
	mb.loadData(ChildMenuDBTable, &result)
	var child []ChildMenu
	for _, v := range result {
		var tmp ChildMenu
		bsBytes, _ := bson.Marshal(v)
		bson.Unmarshal(bsBytes, &tmp)
		child = append(child, tmp)
	}
	return child
}

// load widget list
func (mb *MonitoringBase) LoadWidgetList() []Widget {
	var result []interface{}
	mb.loadData(WidgetListDBTable, &result)
	var wigets []Widget
	for _, v := range result {
		var tmp Widget
		bsBytes, _ := bson.Marshal(v)
		bson.Unmarshal(bsBytes, &tmp)
		wigets = append(wigets, tmp)
	}
	return wigets
}

func (mb *MonitoringBase) LoadWidgetListByID(wgid string) Widget {
	session, err := mb.sessionStart()
	mb.checkError(err)
	c := session.DB(config.Base).C(WidgetListDBTable)
	var result interface{}
	err = c.Find(bson.M{"_id": bson.ObjectIdHex(wgid)}).One(&result)
	var wg Widget
	bsonBytes, _ := bson.Marshal(result)
	bson.Unmarshal(bsonBytes, &wg)
	return wg

}

// Write Device Group list
func (mb *MonitoringBase) WriteDeviceGroup(deviceName string) {
	dev_group := devices.DeviceGroup{bson.NewObjectId(), deviceName}
	mb.insertData(DeviceGroupDBTable, dev_group)
}

// Write Network device list
func (mb *MonitoringBase) WriteNetDev(name string, locate string, ip string, active bool, groupid bson.ObjectId) {
	net_dev := devices.NetDevice{bson.NewObjectId(), name, locate, ip, active, groupid}
	mb.insertData(NetDeviceDBTable, net_dev)
}

// Write OID List
func (mb *MonitoringBase) WriteOidList(name string, oid string, groupid int64, repeat int64) {
	oid_list := devices.OidList{bson.NewObjectId(), name, oid, groupid, repeat}
	mb.insertData(OidListDBTable, oid_list)
}

// Write Menu Group
func (mb *MonitoringBase) WriteMenuGroupList(menuTitle string, pageid string) {
	menuGroup := MenuGroups{bson.NewObjectId(), menuTitle, pageid}
	mb.insertData(MenuGroupDBTable, menuGroup)
}

// write monitoring page
func (mb *MonitoringBase) WriteMonitoringPage(pageName string, pageWg string) {
	page := MonitoringPages{bson.NewObjectId(), pageName, pageWg}
	mb.insertData(MonitorinPagesDBTable, page)
}

// write child menu
func (mb *MonitoringBase) WriteChildMenu(title string, parent string, pageid string) {
	child := ChildMenu{bson.NewObjectId(), title, parent, pageid}
	mb.insertData(ChildMenuDBTable, child)
}

// write widget
func (mb *MonitoringBase) WriteWidgetToBase(wgname string, wgtableName string, wgtype string) {
	wg := Widget{bson.NewObjectId(), wgname, wgtableName, wgtype}
	mb.insertData(WidgetListDBTable, wg)
}

func (mb *MonitoringBase) LoadDataByTableName(table string) interface{} {
	var data interface{}

	switch table {
	case NetDeviceDBTable:
		data = mb.LoadNetDevice()
	case DeviceGroupDBTable:
		data = mb.LoadDeviceGroup()
	default:
		panic("not found table")
	}
	return data
}
