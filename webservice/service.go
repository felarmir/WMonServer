package webservice

import (
	"html/template"
	"log"
	"net/http"
	"path"

	"encoding/json"

	"../datasource"
	"gopkg.in/mgo.v2/bson"
)

var (
	page_template = template.Must(template.ParseFiles(path.Join("webservice/templates", "index.html")))
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

	wg_factory := new(WidgetListCreat)

	//pageData
	pd := PageData{}
	// page scripts
	pd.Tablescripts = true
	pd.ChartScripts = true
	pd.Menu = MenuGenerator(dataLoader.MenuGroupsList()) // left menu
	// widgets
	pd.registerTableWidget(wg_factory.WidgetGenerate(data, 6, "Device group", TABLE_WITH_FORM, "devicegroup").GetWidgetData())
	pd.registerTableWidget(wg_factory.WidgetGenerate(data, 6, "Device group2", EDITABLE_TABLE, "devicegroup").GetWidgetData())

	//pd.registerTableWidget(wg_factory.WidgetGenerate(devList, 12, "Device List", "etable", "netdevice").GetWidgetData())

	err := page_template.ExecuteTemplate(writer, "layout", pd)
	webWerror(err, &writer)
}

//Handler monitor API add
func monitorAPIAdd(writer http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	switch req.Form.Get("datapath") {
	case datasource.DEVICE_GROUP_DBTABLE:
		if len(req.Form.Get("name")) != 0 {
			dataLoader.WriteDeviceGroup(req.Form.Get("name"))
		}
	case datasource.NET_DEVICE_DBTABLE:
		var active bool
		if req.Form.Get("active") == "on" {
			active = true
		} else {
			active = false
		}
		dataLoader.WriteNetDev(req.Form.Get("name"), req.Form.Get("located"), req.Form.Get("ip"), active, bson.ObjectIdHex(req.Form.Get("groupid")))
	case datasource.MENU_GROUP_DBTABLE:
		dataLoader.WriteMenuGroupList(req.Form.Get("title"), req.Form.Get("pageid"))
	case datasource.MONITORING_PAGES_DBTABLE:
		dataLoader.WriteMonitoringPage(req.Form.Get("name"), req.Form.Get("widget"), req.Form.Get("data"))
	case datasource.CHULD_MENU_DBTABLE:
		dataLoader.WriteChildMenu(req.Form.Get("title"), req.Form.Get("parentid"), req.Form.Get("pageid"))

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
	case datasource.DEVICE_GROUP_DBTABLE:
		dataForJSON = dataLoader.LoadDeviceGroup()
	case datasource.NET_DEVICE_DBTABLE:
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
	case datasource.NET_DEVICE_DBTABLE:
		activeBool := false
		if req.Form.Get("Active") == "on" {
			activeBool = true
		}
		dataLoader.UpdateDataRow(datasource.NET_DEVICE_DBTABLE, req.Form.Get("rowID"), bson.M{"name": req.Form.Get("name"), "located": req.Form.Get("located"), "ip": req.Form.Get("ip"), "active": activeBool, "groupid": bson.ObjectIdHex(req.Form.Get("groupid"))})
	case datasource.DEVICE_GROUP_DBTABLE:
		dataLoader.UpdateDataRow(datasource.DEVICE_GROUP_DBTABLE, req.Form.Get("rowID"), bson.M{"name": req.Form.Get("name")})

	default:
		log.Println("not faund table ")
	}
}

// Handler for Page generator Section
func monitoringPages(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/html")

	err := page_template.ExecuteTemplate(writer, "layout", nil)
	webWerror(err, &writer)
}

// handler for settings
func monitorSettings(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/html")

	settingPage := PageData{}
	settingPage.ChartScripts = false
	settingPage.Tablescripts = true
	settingPage.Menu = MenuGenerator(dataLoader.MenuGroupsList())

	wg_factory := new(WidgetListCreat)
	settingPage.registerTableWidget(wg_factory.WidgetGenerate(dataLoader.MenuGroupsList(), 12, "Menu group", TABLE_WITH_FORM, datasource.MENU_GROUP_DBTABLE).GetWidgetData())
	settingPage.registerTableWidget(wg_factory.WidgetGenerate(dataLoader.LoadMonitoringPages(), 12, "Pages", TABLE_WITH_FORM, datasource.MONITORING_PAGES_DBTABLE).GetWidgetData())
	settingPage.registerTableWidget(wg_factory.WidgetGenerate(dataLoader.ChildMenuList(), 12, "Child Menu", TABLE_WITH_FORM, datasource.CHULD_MENU_DBTABLE).GetWidgetData())
	settingPage.registerTableWidget(wg_factory.WidgetGenerate(dataLoader.LoadWidgetList(), 12, "Widget List", TABLE_WITH_FORM, datasource.WIDGET_LIST).GetWidgetData())

	err := page_template.ExecuteTemplate(writer, "layout", settingPage)
	webWerror(err, &writer)
}

// ================== server entry point ===============================
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
