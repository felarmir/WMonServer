package datasource

import "gopkg.in/mgo.v2/bson"

type MonitoringPages struct {
	ID       bson.ObjectId `bson:"_id"`
	Name     string        `bson:"name"`
	WidgetID string        `bson:"widget"`
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
	ID            bson.ObjectId `bson:"_id"`
	Name          string        `bson:"name"`
	DataTableName string        `bson:"datatablename"`
	Groupid		  bson.ObjectId `bson:"datagroupID"`
	WidgetType    string        `bson:"widgettype"`
}
