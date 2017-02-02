package datasource

import "gopkg.in/mgo.v2/bson"

type Page struct {
	ID bson.ObjectId `bson:"_id"`
	PageName string `bson:"pname"`
	WGname string `bson:""wgname`
}