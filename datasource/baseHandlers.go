package datasource

import (
	"../devices"
	"gopkg.in/mgo.v2/bson"
)

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

// Load device by GroupID
func (mb *MonitoringBase) LoadNetDeviceByGroup(groupID bson.ObjectId) []devices.NetDevice {
	var devgroupI []interface{}
	mb.loadDataByCondition(NetDeviceDBTable, &devgroupI, bson.M{"groupid": groupID})
	var netDev []devices.NetDevice

	for _, v := range devgroupI {
		var st devices.NetDevice
		bsonBytes, _ := bson.Marshal(v)
		bson.Unmarshal(bsonBytes, &st)
		netDev = append(netDev, st)
	}
	return netDev
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

func (mb *MonitoringBase) LoadOIDList() []devices.OidList {
	var result []interface{}
	mb.loadData(OidListDBTable, &result)
	var oids []devices.OidList
	for _, v := range result {
		var tmp devices.OidList
		bsBytes, _ := bson.Marshal(v)
		bson.Unmarshal(bsBytes, &tmp)
		oids = append(oids, tmp)
	}
	return oids
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
func (mb *MonitoringBase) WriteOidList(name string, oid string, groupid bson.ObjectId, repeat int64) {
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
func (mb *MonitoringBase) WriteWidgetToBase(wgname string, wgtableName string, groupid bson.ObjectId, wgtype string) {
	wg := Widget{bson.NewObjectId(), wgname, wgtableName, groupid, wgtype}
	mb.insertData(WidgetListDBTable, wg)
}

func (mb *MonitoringBase) LoadDataByTableName(table string, groupID bson.ObjectId) interface{} {
	var data interface{}

	switch table {
	case NetDeviceDBTable:
		data = mb.LoadNetDeviceByGroup(groupID)
	case DeviceGroupDBTable:
		data = mb.LoadDeviceGroup()
	case OidListDBTable:
		data = mb.LoadOIDList()
	default:
		panic("not found table")
	}
	return data
}
