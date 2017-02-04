package datasource

import "gopkg.in/mgo.v2/bson"


type MonitoringPages struct {
	ID     bson.ObjectId `bson:"_id"`
	Name   string        `bson:"name"`
	Widget string        `bson:"widget"`
	Data   string        `bson:"data"`
}

type MenuGroups struct {
	ID                bson.ObjectId `bson:"_id"`
	Title             string        `bson:"menutitle"`
	MonitoringPagesID string        `bson:"pageid"`
}

type ChildMenu struct {
	ID                bson.ObjectId `bson:"_id"`
	Title             string        `bson:"menutitle"`
	MenuGroupID       string        `bson:"menugroupid"`
	MonitoringPagesID string        `bson:"pageid"`
}

type Widget struct {
	ID bson.ObjectId
	Name string
	DataTableName string
	WidgetType string
}