package datasource

//Databese table names
const (
	DeviceGroupDBTable    = "devicegroup"
	NetDeviceDBTable      = "networkdevices"
	OidListDBTable        = "oislist"
	MonitorinPagesDBTable = "monitoringpages"
	MenuGroupDBTable      = "menugroup"
	ChildMenuDBTable      = "childmenu"
	WidgetListDBTable     = "widgetlist"
)

func BaseTablesList() map[string]string {
	baseList := map[string]string{
		"devicegroup":     "Device Group Table",
		"networkdevices":  "Network Device Table",
		"oislist":         "OID Table",
		"monitoringpages": "Monitoring Pages Table",
		"menugroup":       "Menu Group Table",
		"childmenu":       "Child menu Table",
		"widgetlist":      "Widget List Table",
	}
	return baseList
}
