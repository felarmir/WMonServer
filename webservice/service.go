package webservice

import (
	"html/template"
	"log"
	"net/http"
	"path"

	"../datasource"
	"../devices"
)

var (
	page_template = template.Must(template.ParseFiles(path.Join("webservice/templates", "index.html")))
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

	base := datasource.MonitoringBase{}
	data := base.LoadDeviceGroup()

	wg_factory := new(WidgetListCreat)

	//pageData
	pd := PageData{}
	pd.registerTableWidget(wg_factory.WidgetGenerate(data, 6, "Device group", "table").GetWidgetData())
	pd.registerTableWidget(wg_factory.WidgetGenerate(data, 6, "Device group2", "table").GetWidgetData())

	pd.registerFormWidget(wg_factory.WidgetGenerate(devices.NetDev{}, 12, "Devise add", "form").GetWidgetData())

	err := page_template.ExecuteTemplate(writer, "layout", pd)
	webWerror(err, &writer)
}

//Handler for  monitoring
func monitorMonitorHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/html")
	err := page_template.ExecuteTemplate(writer, "layout", nil)
	webWerror(err, &writer)
}

// Handler for settings Section
func monitoringManagingHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/html")
	err := page_template.ExecuteTemplate(writer, "layout", nil)
	webWerror(err, &writer)
}

func WebServer() {
	fs := http.FileServer(http.Dir("./webservice/public/static")) // static files real path
	http.Handle("/static/", http.StripPrefix("/static/", fs))     // static files path

	http.HandleFunc("/", monitorIndexHandler)
	http.HandleFunc("/monitor", monitorMonitorHandler)
	http.HandleFunc("/settings", monitoringManagingHandler)

	log.Println("Server start ...")
	http.ListenAndServe(":8000", nil)
}
