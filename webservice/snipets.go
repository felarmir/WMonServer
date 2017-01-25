package webservice

import (
	"html/template"
	"reflect"
	"strconv"
	"log"
)

type Widgets interface {
	WidgetGenerate(data interface{}, widgetSize int64, widgetTitle string) Widget
	registerWidget(widget Widget)
}
type Widget interface {
	GetWidgetData() template.HTML
}

type WidgetListCreat struct {
	widgetArray []*Widget
}

func (self *WidgetListCreat) registerWidget(widget Widget) {
	self.widgetArray = append(self.widgetArray, &widget)
}

type TableWidget struct {
	Tabledata template.HTML
}

func (self *TableWidget) GetWidgetData() template.HTML {
	return self.Tabledata
}


// Table widget head part
// tableSize: 3 ... 12; tableName Table name
func getHeaderDIV(tableSize int64, tableName string) string {
	snp := "<div class=\"col-md-" + strconv.FormatInt(tableSize, 10) + "\"> <div class=\"content-box-large\"> " +
		"<div class=\"panel-heading\"> <div class=\"panel-title\">" + tableName + "</div> " +
		"<div class=\"panel-options\"> <a href=\"#\" data-rel=\"collapse\"><i class=\"glyphicon glyphicon-refresh\"></i></a> " +
		"<a href=\"#\" data-rel=\"reload\"><i class=\"glyphicon glyphicon-cog\"></i></a> " +
		"</div> </div> <div class=\"panel-body\"> <table class=\"table table-bordered\"> "
	return snp
}

// Parse interface and get column names from fields name
// Generate Table Header
func generateTableHeader(data interface{}) string {
	topen := "<thead> <tr>"
	tclose := "</tr></thead>"
	var header string

	preF := reflect.TypeOf(data).Elem()
	for i := 0; i < preF.NumField(); i++ {
		header = header + "<th>" + preF.Field(i).Name + "</th>"
	}
	return topen + header + tclose
}

//Generate Table content
// generate table with data
func generateTableData(data interface{}) string {
	result := "<tbody>"

	s := reflect.ValueOf(data)
	if s.Kind() != reflect.Slice {
		panic("Data not a slice type")
	}
	vdata := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		vdata[i] = s.Index(i).Interface()
	}

	for _, v := range vdata {
		tmp := "<tr>"
		pre := reflect.ValueOf(v)
		for i := 0; i < pre.NumField(); i++ {
			var value string
			if r, ok := pre.Field(i).Interface().(string); ok {
				value = r
			}
			if r, ok := pre.Field(i).Interface().(int64); ok {
				value = strconv.FormatInt(r, 10)
			}
			tmp = tmp + "<td>" + value + "</td>"
		}
		result = result + tmp + "</tr>"
	}

	return result + "</tbody>"
}

//Function for generate table widget
// data from base; tableSize: 3 ... 12
func tableWidgetGenerator(data interface{}, tableSize int64, tableName string) template.HTML {
	var result string
	result = getHeaderDIV(tableSize, tableName) + generateTableHeader(data) + generateTableData(data)
	result += "</table> </div> </div> </div>"
	return template.HTML([]byte(result))
}

func (self *WidgetListCreat) WidgetGenerate(data interface{}, widgetSize int64, widgetTitle string, widgetType string) Widget {
	var wg Widget

	switch widgetType {
	case "table":
		wg = &TableWidget{tableWidgetGenerator(data, widgetSize, widgetTitle)}
	default:
		log.Fatalln("Unknown Error")
	}

	self.registerWidget(wg)
	return wg
}