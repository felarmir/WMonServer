package webservice

import (
	"net/http"
	"html/template"
	"path"
	"log"
)

var  (
	page_template = template.Must(template.ParseFiles(path.Join("templates", "index.html")))
)

func monitorIndexHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html")
	err := page_template.ExecuteTemplate(res, "layout", nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(res, http.StatusText(500), 500)
	}
}

func WebServer() {
	fs := http.FileServer(http.Dir("./public/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", monitorIndexHandler)
	log.Println("Server start ...")
	http.ListenAndServe(":8000", nil)
}
