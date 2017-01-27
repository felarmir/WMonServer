package webservice

import (
	"html/template"
	"log"
	"reflect"
	"strconv"
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

//Function for generate table widget. Parsing interface data
// data from base; tableSize: 3 ... 12
func tableWidgetGenerator(data interface{}, tableSize int64, tableName string) template.HTML {
	result := "<div class=\"col-md-" + strconv.FormatInt(tableSize, 10) + "\"><div class=\"panel panel-default\">" +
		"<div class=\"panel-heading\"><h3 class=\"panel-title\">" + tableName + "</h3>" +
		"</div><div class=\"panel-body\"><div class=\"row\"><div class=\"col-md-12 col-sm-12 col-xs-12\">" +
		"<table class=\"table\"><thead><tr>"

	// add table header
	preF := reflect.TypeOf(data).Elem()
	for i := 0; i < preF.NumField(); i++ {
		result = result + "<th>" + preF.Field(i).Name + "</th>"
	}
	result = result + "</tr></thead><tbody>"

	// add table content
	dataSlice := reflect.ValueOf(data)
	if dataSlice.Kind() != reflect.Slice {
		panic("Data not a slice type")
	}
	vdata := make([]interface{}, dataSlice.Len())
	for i := 0; i < dataSlice.Len(); i++ {
		vdata[i] = dataSlice.Index(i).Interface()
	}

	for _, v := range vdata {
		tmp := "<tr class=\"active\">"
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

	//result = result + generateTableData(data)
	result += "</tbody></table> </div> </div> </div></div> </div>"
	return template.HTML([]byte(result))
}

// Function Editable Table Generate
//data is object from stuct
func editableTableWidgetGenerate(data interface{}, widgetSize int64, widgetTitle string) template.HTML {
	editTable :=  "<div class=\"row\"><div class=\"col-sm-8\"><h4 class=\"pull-left page-title\">Editable Table</h4></div></div>"+
	"<div class=\"col-sm-8\"><div class=\"panel\"><div class=\"panel-body\"><div class=\"row\"><div class=\"col-sm-6\"><div class=\"m-b-30\">"+
	"<button id=\"addToTable\" class=\"btn btn-primary waves-effect waves-light\">Add <i class=\"fa fa-plus\"></i></button>"+
	"</div> </div> </div> <table class=\"table table-bordered table-striped\" id=\"datatable-editable\"><thead><tr>"

	// add table header
		preF := reflect.TypeOf(data).Elem()
	for i := 0; i < preF.NumField(); i++ {
		editTable = editTable + "<th>" + preF.Field(i).Name + "</th>"
	}
	editTable += "<th>Actions</th></tr></thead><tbody>"
	// add table content
	dataSlice := reflect.ValueOf(data)
	if dataSlice.Kind() != reflect.Slice {
		panic("Data not a slice type")
	}
	vdata := make([]interface{}, dataSlice.Len())
	for i := 0; i < dataSlice.Len(); i++ {
		vdata[i] = dataSlice.Index(i).Interface()
	}

	for _, v := range vdata {
		tmp := "<tr class=\"gradeX\">"
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
		editTable = editTable + tmp + "<td class=\"actions\">"+
			"<a href=\"#\" class=\"hidden on-editing save-row\"><i class=\"fa fa-save\"></i></a>"+
			"<a href=\"#\" class=\"hidden on-editing cancel-row\"><i class=\"fa fa-times\"></i></a>"+
			"<a href=\"#\" class=\"on-default edit-row\"><i class=\"fa fa-pencil\"></i></a>"+
			"<a href=\"#\" class=\"on-default remove-row\"><i class=\"fa fa-trash-o\"></i></a>"+
			"</td></tr>"
	}
	editTable += "</tbody></table></div></div> </div>"

	return template.HTML([]byte(editTable))
}


// Function for input form generate
//data is object from stuct
func formWidgetGenerator(data interface{}, widgetSize int64, widgetTitle string) template.HTML {
	form := "<div class=\"col-md-" + strconv.FormatInt(widgetSize, 10) + "\"><div class=\"panel panel-default\">" +
		"<div class=\"panel-heading\"><h3 class=\"panel-title\">" + widgetTitle + "</h3></div>" +
		"<div class=\"panel-body\"><form class=\"form-horizontal\" role=\"form\">"

	preF := reflect.TypeOf(data)
	for i := 0; i < preF.NumField(); i++ {
		form = form + "<div class=\"form-group\"><label class=\"col-md-2 control-label\">" + preF.Field(i).Name + "</label>" +
			"<div class=\"col-md-10\"><input type=\"text\" name=\"" + preF.Field(i).Name + "\" class=\"form-control\" value=\"\">" +
			"</div></div>"
	}

	form = form + "<div class=\"form-group m-b-0\"><div class=\"col-sm-offset-3 col-sm-9\">" +
		"<button type=\"submit\" class=\"btn btn-info waves-effect waves-light\">Sign in</button>" +
		"</div> </div></form> </div></div></div>"
	return template.HTML([]byte(form))
}

func (self *WidgetListCreat) WidgetGenerate(data interface{}, widgetSize int64, widgetTitle string, widgetType string) Widget {
	var wg Widget

	switch widgetType {
	case "table":
		wg = &ReadyWidget{tableWidgetGenerator(data, widgetSize, widgetTitle)}
	case "etable":
		wg = &ReadyWidget{editableTableWidgetGenerate(data, widgetSize, widgetTitle)}
	case "form":
		wg = &ReadyWidget{formWidgetGenerator(data, widgetSize, widgetTitle)}

	default:
		log.Fatalln("Unknown Error")
	}

	self.registerWidget(wg)
	return wg
}
