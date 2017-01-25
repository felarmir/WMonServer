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

type ReadyWidget struct {
	Tabledata template.HTML
}

func (self *ReadyWidget) GetWidgetData() template.HTML {
	return self.Tabledata
}


// Table widget head part
// tableSize: 3 ... 12; tableName Table name
func getHeaderDIV(tableSize int64, tableName string) string {
	 snp := "<div class=\"col-md-"+ strconv.FormatInt(tableSize, 10) +" col-sm-6 col-xs-12\"><div class=\"x_panel\">"+
		"<div class=\"x_title\"><h2>"+ tableName +"</h2>"+
		"<ul class=\"nav navbar-right panel_toolbox\"><li><a class=\"collapse-link\">"+
		"<i class=\"fa fa-chevron-up\"></i></a></li>"+
		 "</ul><div class=\"clearfix\"></div></div><div class=\"x_content\"><table class=\"table table-striped\">"

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


// Function for input form generate
func formWidgetGenerator(data interface{}, widgetSize int64, widgetTitle string) template.HTML {
 	form := "<div class=\"col-md-"+strconv.FormatInt(widgetSize, 10)+" col-sm-12 col-xs-12\"> <div class=\"x_panel\"> <div class=\"x_title\">"+
	"<h2>"+widgetTitle+"</h2> <ul class=\"nav navbar-right panel_toolbox\">"+
	"<li><a class=\"collapse-link\"><i class=\"fa fa-chevron-up\"></i></a></li> </ul> <div class=\"clearfix\"></div>" +
	"</div> <div class=\"x_content\"><br><form id=\"demo-form2\" data-parsley-validate=\"\" class=\"form-horizontal form-label-left\" novalidate=\"\">"

	preF := reflect.TypeOf(data)
	for i := 0; i < preF.NumField(); i++ {
		form = form + "<div class=\"form-group\"><label class=\"control-label col-md-3 col-sm-3 col-xs-12\">"+preF.Field(i).Name+" <span class=\"required\">*</span></label>"+
		"<div class=\"col-md-6 col-sm-6 col-xs-12\">"+
		"<input type=\"text\" required=\"required\" name=\""+preF.Field(i).Name+"\" class=\"form-control col-md-7 col-xs-12\"></div></div>"
	}

	form = form + "<div class=\"ln_solid\"></div> <div class=\"form-group\"><div class=\"col-md-6 col-sm-6 col-xs-12 col-md-offset-3\">"+
		"<button class=\"btn btn-primary\" type=\"button\">Cancel</button><button class=\"btn btn-primary\" type=\"reset\">Reset</button>"+
		"<button type=\"submit\" class=\"btn btn-success\">Submit</button></div> </div> </form> </div> </div> </div>"
	return template.HTML([]byte(form))
}


func (self *WidgetListCreat) WidgetGenerate(data interface{}, widgetSize int64, widgetTitle string, widgetType string) Widget {
	var wg Widget

	switch widgetType {
	case "table":
		wg = &ReadyWidget{tableWidgetGenerator(data, widgetSize, widgetTitle)}
	case "form":
		wg = &ReadyWidget{formWidgetGenerator(data, widgetSize, widgetTitle)}

	default:
		log.Fatalln("Unknown Error")
	}

	self.registerWidget(wg)
	return wg
}