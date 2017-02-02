package datasource

import "gopkg.in/mgo.v2/bson"

type Page struct {
	ID       bson.ObjectId `bson:"_id"`
	PageName string        `bson:"pname"`
	WGname   string        `bson:""wgname`
}

type SideMenuList struct {
	ID        bson.ObjectId  `bson:"_id"`
	MenuTitle string         `bson:"menutitle"`
	Pageid    string         `bson:"pageid"`
	ChildNode []SideMenuList `bson:"childnode"`
}
