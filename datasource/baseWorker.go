package datasource

import (
	"fmt"

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

// Insert interface Data to table
func (mb *MonitoringBase) insertData(table string, data interface{}) {
	session, err := mb.sessionStart()
	mb.checkError(err)
	c := session.DB(config.Base).C(table)
	err = c.Insert(&data)
	mb.checkError(err)
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

// Load data by table name, condition and return intrface
func (mb *MonitoringBase) loadDataByCondition(table string, data *[]interface{}, condition bson.M) {
	session, err := mb.sessionStart()
	mb.checkError(err)
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(config.Base).C(table)
	var result []interface{}
	err = c.Find(condition).All(&result)
	mb.checkError(err)
	*data = result
}
