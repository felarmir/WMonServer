package webservice

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strconv"

	"encoding/json"

	"../datasource"
	"gopkg.in/mgo.v2/bson"
)

var (
	pageTemplate = template.Must(template.ParseFiles(path.Join("webservice/templates", "index.html")))
)

var (
	dataLoader datasource.MonitoringBase
)

func webWerror(err error, res *http.ResponseWriter) {
	if err != nil {
		log.Println(err.Error())
		http.Error(*res, http.StatusText(500), 500)
	}
}

// index page handler
func monitorIndexHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/html")

	data := dataLoader.LoadDeviceGroup()

	//devList := dataLoader.LoadNetDevice()

	wgfactory := new(WidgetListCreat)

	//pageData
	pd := PageData{}
	// page scripts
	pd.Tablescripts = true
	pd.ChartScripts = true
	pd.Menu = MenuGenerator(dataLoader.MenuGroupsList()) // left menu
	// widgets
	pd.registerTableWidget(wgfactory.WidgetGenerate(data, 6, "Device group", TableWithForm, "devicegroup").GetWidgetData())
	pd.registerTableWidget(wgfactory.WidgetGenerate(data, 6, "Device group2", EditbleTable, "devicegroup").GetWidgetData())

	//pd.registerTableWidget(wgfactory.WidgetGenerate(devList, 12, "Device List", "etable", "netdevice").GetWidgetData())

	err := pageTemplate.ExecuteTemplate(writer, "layout", pd)
	webWerror(err, &writer)
}

//Handler monitor API add
func monitorAPIAdd(writer http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	switch req.Form.Get("datapath") {
	case datasource.DeviceGroupDBTable:
		if len(req.Form.Get("name")) != 0 {
			dataLoader.WriteDeviceGroup(req.Form.Get("name"))
		}

	case datasource.NetDeviceDBTable:
		var active bool
		if req.Form.Get("active") == "on" {
			active = true
		} else {
			active = false
		}
		dataLoader.WriteNetDev(req.Form.Get("name"), req.Form.Get("located"), req.Form.Get("ip"), active, bson.ObjectIdHex(req.Form.Get("groupid")))

	case datasource.MenuGroupDBTable:
		dataLoader.WriteMenuGroupList(req.Form.Get("title"), req.Form.Get("monitoringpagesid"))

	case datasource.MonitorinPagesDBTable:
		dataLoader.WriteMonitoringPage(req.Form.Get("name"), req.Form.Get("widgetid"))

	case datasource.ChildMenuDBTable:
		dataLoader.WriteChildMenu(req.Form.Get("title"), req.Form.Get("menugroupid"), req.Form.Get("monitoringpagesid"))

	case datasource.WidgetListDBTable:
		sv := req.Form.Get("widgettype")
		wgType, _ := strconv.ParseInt(sv, 10, 64)
		dataLoader.WriteWidgetToBase(req.Form.Get("name"), req.Form.Get("datatablename"), WidgetTypeMap()[wgType])
	default:
		log.Panicln("Undefine table")
	}

	http.Redirect(writer, req, "/", 301)
}

// Header for Api get json
func monitorAPIGetJSON(writer http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	var dataForJSON interface{}

	switch req.Form.Get("name") {
	case datasource.DeviceGroupDBTable:
		dataForJSON = dataLoader.LoadDeviceGroup()
	case datasource.NetDeviceDBTable:
		dataForJSON = dataLoader.LoadNetDevice()
	default:
		log.Println("Error load Data")
	}

	jsn, _ := json.Marshal(struct {
		Result string      `json:"Result"`
		AaData interface{} `json:"Records"`
	}{"OK", dataForJSON})

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(jsn)
}

// Handler for Api delete Row
func monitorAPIDeleteRow(writer http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	tableName := req.Form.Get("datapath")
	rowID := req.Form.Get("rowID")
	dataLoader.DeleteDataRow(tableName, rowID)
}

//Handler for Api Update Row
func monitorAPIUpdateRow(writer http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	switch req.Form.Get("datapath") {
	case datasource.NetDeviceDBTable:
		activeBool := false
		if req.Form.Get("Active") == "on" {
			activeBool = true
		}
		dataLoader.UpdateDataRow(datasource.NetDeviceDBTable, req.Form.Get("rowID"), bson.M{"name": req.Form.Get("name"), "located": req.Form.Get("located"), "ip": req.Form.Get("ip"), "active": activeBool, "groupid": bson.ObjectIdHex(req.Form.Get("groupid"))})
	case datasource.DeviceGroupDBTable:
		dataLoader.UpdateDataRow(datasource.DeviceGroupDBTable, req.Form.Get("rowID"), bson.M{"name": req.Form.Get("name")})

	default:
		log.Println("not faund table ")
	}
}

// Handler for Page generator Section
func monitoringPages(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/html")
	req.ParseForm()
	pageID := req.Form.Get("id")

	pageData := dataLoader.LoadMonitoringPage(pageID)
	wg := dataLoader.LoadWidgetListByID(pageData.WidgetID)

	dynPage := PageData{}
	dynPage.ChartScripts = false
	dynPage.Tablescripts = true
	dynPage.Menu = MenuGenerator(dataLoader.MenuGroupsList())
	wgfactory := new(WidgetListCreat)

	dynPage.registerTableWidget(wgfactory.WidgetGenerate(dataLoader.LoadDataByTableName(wg.DataTableName), 12, wg.Name, wg.WidgetType, wg.DataTableName).GetWidgetData())

	err := pageTemplate.ExecuteTemplate(writer, "layout", dynPage)

	webWerror(err, &writer)
}

// handler for settings
func monitorSettings(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/html")

	settingPage := PageData{}
	settingPage.ChartScripts = false
	settingPage.Tablescripts = true
	settingPage.Menu = MenuGenerator(dataLoader.MenuGroupsList())

	wgfactory := new(WidgetListCreat)
	settingPage.registerTableWidget(wgfactory.WidgetGenerate(dataLoader.MenuGroupsList(), 12, "Menu group", TableWithForm, datasource.MenuGroupDBTable).GetWidgetData())
	settingPage.registerTableWidget(wgfactory.WidgetGenerate(dataLoader.LoadMonitoringPages(), 12, "Pages", TableWithForm, datasource.MonitorinPagesDBTable).GetWidgetData())
	settingPage.registerTableWidget(wgfactory.WidgetGenerate(dataLoader.ChildMenuList(), 12, "Child Menu", TableWithForm, datasource.ChildMenuDBTable).GetWidgetData())
	settingPage.registerTableWidget(wgfactory.WidgetGenerate(dataLoader.LoadWidgetList(), 12, "Widget List", TableWithForm, datasource.WidgetListDBTable).GetWidgetData())

	err := pageTemplate.ExecuteTemplate(writer, "layout", settingPage)
	webWerror(err, &writer)
}

// WebServer entry point
func WebServer() {
	fs := http.FileServer(http.Dir("./webservice/public/static")) // static files real path
	http.Handle("/static/", http.StripPrefix("/static/", fs))     // static files path

	dataLoader = datasource.MonitoringBase{}

	http.HandleFunc("/", monitorIndexHandler)
	http.HandleFunc("/page", monitoringPages)

	http.HandleFunc("/settings", monitorSettings)

	http.HandleFunc("/api/add", monitorAPIAdd)
	http.HandleFunc("/api/get", monitorAPIGetJSON)
	http.HandleFunc("/api/del", monitorAPIDeleteRow)
	http.HandleFunc("/api/update", monitorAPIUpdateRow)

	log.Println("Server start ...")
	http.ListenAndServe(":8000", nil)
}
