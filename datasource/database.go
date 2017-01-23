package datasource

import (
	"fmt"

	"../devices"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var netnode []devices.NetNode
var snmpTable []devices.SNMPTemplate

func (self *MonitoringBase) CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func (self *MonitoringBase) ConnectSession(user string, pass string, host string, port string) (*mgo.Session, error) {
	if len(port) == 0 {
		port = "27017"
	}
	url := "mongodb://" + user + ":" + pass + "@" + host + ":" + port
	return mgo.Dial(url)
}

func (self *MonitoringBase) ConnectSessionWithDefaultPort(user string, pass string, host string) (*mgo.Session, error) {
	return self.ConnectSession(user, pass, host, "")
}

func (self *MonitoringBase) LoadData(table string) {
	session, err := self.ConnectSessionWithDefaultPort("name", "pass", "ip")
	self.CheckError(err)
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("monitoring").C(table)
	var result []interface{}
	err = c.Find(bson.M{}).All(&result)
	if table == "snmptpl" {
		fmt.Println(result)
	}
	/*
		switch table {
		case "snmptpl":
			err = c.Find(bson.M{}).All(&snmpTable)
			fmt.Println(snmpTable)
		case "netnode":
			err = c.Find(bson.M{}).All(&netnode)
			fmt.Println(snmpTable)
		}
	*/
}

//func (self *MonitoringBase) InsertData(table string, )

/*
func MongoTest() {

	c := session.DB("monitoring").C("people")
	err = c.Insert(&Person{"Ale", "+55 53 8116 9639"},
		&Person{"Cla", "+55 53 8402 8510"})
	if err != nil {
		log.Fatal(err)
	}

}
*/
