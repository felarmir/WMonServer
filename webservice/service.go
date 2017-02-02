package webservice

import (
	"html/template"
	"log"
	"net/http"
	"path"

	"../datasource"
	"encoding/json"
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
	pd.registerTableWidget(wg_factory.WidgetGenerate(data, 6, "Device group", "tablein", "devicegroup").GetWidgetData())
	pd.registerTableWidget(wg_factory.WidgetGenerate(data, 6, "Device group2", "etable", "devicegroup").GetWidgetData())

	//pd.registerTableWidget(wg_factory.WidgetGenerate(devList, 12, "Device List", "etable", "netdevice").GetWidgetData())

	err := page_template.ExecuteTemplate(writer, "layout", pd)
	webWerror(err, &writer)
}

//Handler monitor API add
func monitorAPIAdd(writer http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	switch req.Form.Get("datapath") {
	case "devicegroup":
		if len(req.Form.Get("name")) != 0 {
			dataLoader.WriteDeviceGroup(req.Form.Get("name"))
		}
	case "netdevice":
		var active bool
		if req.Form.Get("active") == "on" {
			active = true
		} else {
			active = false
		}
		dataLoader.WriteNetDev(req.Form.Get("name"), req.Form.Get("located"), req.Form.Get("ip"), active, bson.ObjectIdHex(req.Form.Get("groupid")))
	case "menugroup":
		dataLoader.WriteMenuGroupList(req.Form.Get("title"), req.Form.Get("pageid"), "menugroup")
	case "pages":
		dataLoader.WriteMonitoringPage(req.Form.Get("name"), req.Form.Get("widget"), req.Form.Get("data"), "pages")
	case "childmenu":
		dataLoader.WriteChildMenu(req.Form.Get("title"), req.Form.Get("parentid"), req.Form.Get("pageid"), "childmenu")

	default:
		log.Panicln("Undefine table")
	}

	http.Redirect(writer, req, "/", 301)
}

// Header for Api get json
func monitoringAPIGetJSON(writer http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	var dataForJSON interface{}

	switch req.Form.Get("name") {
	case "devicegroup":
		dataForJSON = dataLoader.LoadDeviceGroup()
	case "netdevice":
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
func monitoringAPIDeleteRow(writer http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	tableName := req.Form.Get("datapath")
	rowID := req.Form.Get("rowID")
	dataLoader.DeleteDataRow(tableName, rowID)
}

//Handler for Api Update Row
func monitoringAPIUpdateRow(writer http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	switch req.Form.Get("datapath") {
	case "netdevice":
		activeBool := false
		if req.Form.Get("Active") == "on" {
			activeBool = true
		}
		dataLoader.UpdateDataRow("netdevice", req.Form.Get("rowID"), bson.M{"name": req.Form.Get("name"), "located": req.Form.Get("located"), "ip": req.Form.Get("ip"), "active": activeBool, "groupid": bson.ObjectIdHex(req.Form.Get("groupid"))})
	case "devicegroup":
		dataLoader.UpdateDataRow("devicegroup", req.Form.Get("rowID"), bson.M{"name": req.Form.Get("name")})

	default:
		log.Println("not faund table ")
	}
}

//Handler for  monitoring
func monitorMonitorHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/html")
	err := page_template.ExecuteTemplate(writer, "layout", nil)
	webWerror(err, &writer)
}

// Handler for Page generator Section
func monitoringPages(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/html")

	err := page_template.ExecuteTemplate(writer, "layout", nil)
	webWerror(err, &writer)
}

// handler for settings
func monitoringSettings(writer http.ResponseWriter, req *http.Request){
	writer.Header().Set("Content-Type", "text/html")

	settingPage := PageData{}
	settingPage.ChartScripts = false
	settingPage.Tablescripts = true
	settingPage.Menu = MenuGenerator(dataLoader.MenuGroupsList())

	wg_factory := new(WidgetListCreat)
	settingPage.registerTableWidget(wg_factory.WidgetGenerate(dataLoader.MenuGroupsList(), 12, "Menu group", "tablein", "menugroup").GetWidgetData())
	settingPage.registerTableWidget(wg_factory.WidgetGenerate(dataLoader.LoadMonitoringPages(), 12, "Pages", "tablein", "pages").GetWidgetData())
	settingPage.registerTableWidget(wg_factory.WidgetGenerate(dataLoader.ChildMenuList(), 12, "Child Menu", "tablein", "childmenu").GetWidgetData())


	err := page_template.ExecuteTemplate(writer, "layout", settingPage)
	webWerror(err, &writer)
}


// ================== server entry point ===============================
func WebServer() {
	fs := http.FileServer(http.Dir("./webservice/public/static")) // static files real path
	http.Handle("/static/", http.StripPrefix("/static/", fs))     // static files path

	dataLoader = datasource.MonitoringBase{}

	http.HandleFunc("/", monitorIndexHandler)
	http.HandleFunc("/monitor", monitorMonitorHandler)
	http.HandleFunc("/page", monitoringPages)

	http.HandleFunc("/settings", monitoringSettings)

	http.HandleFunc("/api/add", monitorAPIAdd)
	http.HandleFunc("/api/get", monitoringAPIGetJSON)
	http.HandleFunc("/api/del", monitoringAPIDeleteRow)
	http.HandleFunc("/api/update", monitoringAPIUpdateRow)

	log.Println("Server start ...")
	http.ListenAndServe(":8000", nil)
}
