package datasource

//Databese table names
const (
	DeviceGroupDBTable    = "devicegroup"
	NetDeviceDBTable      = "networkdevices"
	OidListDBTable        = "oidlist"
	MonitorinPagesDBTable = "monitoringpages"
	MenuGroupDBTable      = "menugroup"
	ChildMenuDBTable      = "childmenu"
	WidgetListDBTable     = "widgetlist"
)

func BaseTablesList() map[string]string {
	baseList := map[string]string{
		DeviceGroupDBTable:     "Device Group Table",
		NetDeviceDBTable:  "Network Device Table",
		OidListDBTable:         "OID Table",
		MonitorinPagesDBTable: "Monitoring Pages Table",
		MenuGroupDBTable:       "Menu Group Table",
		ChildMenuDBTable:       "Child menu Table",
		WidgetListDBTable:      "Widget List Table",
	}
	return baseList
}
