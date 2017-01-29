package devices

import (
	"gopkg.in/mgo.v2/bson"
)

type DeviceGroup struct {
	ID   bson.ObjectId `bson:"_id,omitempty"`
	Name string        `bson:"name"`
}

type NetDev struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	Name    string        `bson:"name"`
	Located string        `bson:"located"`
	IP      string        `bson:"ip"`
	Active  bool          `bson:"active"`
	Groupid bson.ObjectId `bson:"groupid"`
}

type OidList struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	Name    string        `bson:"name"`
	Oid     string        `bson:"oid"`
	Groupid int64         `bson:"groupid"`
	Repeat  int64         `bson:"repeat"`
}
