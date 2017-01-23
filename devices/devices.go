package devices

type DeviceGroup struct {
	ID   int64
	Name string
}

type NetNode struct {
	ID       int32
	Name     string
	Location string
	IP       string
	Active   bool
	GroupID  int64
}

type SNMPTemplate struct {
	GroupID int64
	OID     map[string]map[string]string
}
