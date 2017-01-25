package webservice

import (
	"net/http"
	"html/template"
	"path"
	"log"
	"../datasource"
)

var  (
	page_template = template.Must(template.ParseFiles(path.Join("webservice/templates", "index.html")))
)

func webWerror(err error, res *http.ResponseWriter) {
	if err != nil {
		log.Println(err.Error())
		http.Error(*res, http.StatusText(500), 500)
	}
}


func monitorIndexHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/html")

	base := datasource.MonitoringBase{}
	data := base.LoadDeviceGroup()

	err := page_template.ExecuteTemplate(writer, "layout", TableWidgetGenerator(data, 3, "Device Group"))
	webWerror(err, &writer)
}

func monitorCalendarHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/html")
	err := page_template.ExecuteTemplate(writer, "layout", nil)
	webWerror(err, &writer)
}

func monitoringManagingHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/html")

}

func WebServer() {
	fs := http.FileServer(http.Dir("./webservice/public/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", monitorIndexHandler)
	http.HandleFunc("/calendar", monitorCalendarHandler)
	log.Println("Server start ...")
	http.ListenAndServe(":8000", nil)
}
