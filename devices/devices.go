package devices



type DeviceGroup struct {
	ID   int64  `bson:"id"`
	Name string `bson:"name"`
}

type NetDev struct {
	ID      int64  `bson:"id"`
	Name    string `bson:"name"`
	Located string `bson:"located"`
	IP      string `bson:"ip"`
	Active  bool   `bson:"active"`
	Groupid int64  `bson:"groupid"`
}

type OidList struct {
	Name    string `bson:"name"`
	Oid     string `bson:"oid"`
	Groupid int64  `bson:"groupid"`
	Repeat  int64  `bson:"repeat"`
}
