package webservice

import (
	"../datasource"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"log"
	"reflect"
	"strconv"
	"strings"
)

var dataSource datasource.MonitoringBase

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
func tableGenerator(data interface{}) string {
	tb := "<table class=\"table\"><thead><tr>"

	// add table header
	preF := reflect.TypeOf(data).Elem()
	for i := 0; i < preF.NumField(); i++ {
		if preF.Field(i).Name != "ID" {
			tb = tb + "<th>" + preF.Field(i).Name + "</th>"
		}
	}
	tb += "</tr></thead><tbody>"

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
			if i != 0 {
				var value string
				if r, ok := pre.Field(i).Interface().(string); ok {
					value = r
				}
				if r, ok := pre.Field(i).Interface().(int64); ok {
					value = strconv.FormatInt(r, 10)
				}
				if r, ok := pre.Field(i).Interface().(bool); ok {
					value = strconv.FormatBool(r)
				}
				if r, ok := pre.Field(i).Interface().(bson.ObjectId); ok {
					if preF.Field(i).Name == "Groupid" {
						for _, rowS := range dataSource.LoadDeviceGroup() {
							if rowS.ID == r {
								value = rowS.Name
							}
						}
					}
				}

				tmp = tmp + "<td>" + value + "</td>"
			}
		}
		tb = tb + tmp + "</tr>"
	}
	tb += "</tbody></table>"
	return tb
}

func tableWidgetGenerator(data interface{}, tableSize int64, tableName string) template.HTML {
	result := "<div class=\"col-md-" + strconv.FormatInt(tableSize, 10) + "\"><div class=\"panel panel-default\">" +
		"<div class=\"panel-heading\"><h3 class=\"panel-title\">" + tableName + "</h3>" +
		"</div><div class=\"panel-body\"><div class=\"row\"><div class=\"col-md-12 col-sm-12 col-xs-12\">"

	result += tableGenerator(data)
	result += "</div> </div> </div></div> </div>"
	return template.HTML([]byte(result))
}

// Function for input form generate
//data is object from stuct

func formGenerator(data interface{}, datatable string) string {
	frm := "<form class=\"form-horizontal\" action=\"/api/add\" role=\"form\">" +
		"<input type=\"hidden\" name=\"datapath\" value=\"" + datatable + "\">"

	dataSlice := reflect.ValueOf(data)
	//var element interface{}

	var element reflect.Type
	if dataSlice.Kind() != reflect.Slice {
		panic("Data not a slice type")
	}

	if dataSlice.Len() == 0 {
		element = reflect.TypeOf(data).Elem()
	} else {
		element = reflect.TypeOf(dataSlice.Index(0).Interface())
	}

	preF := element //reflect.TypeOf(data)

	for i := 0; i < preF.NumField(); i++ {
		if preF.Field(i).Name != "ID" {
			frm = frm + "<div class=\"form-group\"><label class=\"col-md-2 control-label\">" + preF.Field(i).Name + "</label><div class=\"col-md-10\">" +
				getFieldType(preF.Field(i).Name) + "</div></div>"
		}

	}
	frm += "<div class=\"form-group m-b-0\"><div class=\"col-sm-offset-9 col-sm-9\">" +
		"<button type=\"submit\" class=\"btn btn-info waves-effect waves-light\">Submit</button>" +
		"</div> </div></form>"
	return frm
}

// Generate Field by struct row
func getFieldType(ftype string) string {
	var result string

	//<select>
	selectStart := func(selectName string) string {
		return "<select class=\"form-control\" name=\"" + selectName + "\">"
	}
	// <option>
	optionRow := func(dataValue string, dataName string) string {
		return "<option value=\"" + dataValue + "\">" + dataName + "</option>"
	}

	switch ftype {
	case "Groupid":
		result += selectStart(strings.ToLower(ftype))
		for _, v := range dataSource.LoadDeviceGroup() {
			result += optionRow(v.ID.Hex(), v.Name)
		}
		result += "</select>"

	case "MenuGroupID":
		result += selectStart(strings.ToLower(ftype))
		for _, v := range dataSource.MenuGroupsList() {
			result += optionRow(v.ID.Hex(), v.Title)
		}
		result += "</select>"

	case "MonitoringPagesID":
		result += selectStart(strings.ToLower(ftype))
		for _, v := range dataSource.LoadMonitoringPages() {
			result += optionRow(v.ID.Hex(), v.Name)
		}
		result += "</select>"

	case "Active":
		result = "<input id=\"checkbox2\" name=\"" + strings.ToLower(ftype) + "\" type=\"checkbox\">"

	default:
		result = "<input type=\"text\" name=\"" + strings.ToLower(ftype) + "\" class=\"form-control\" value=\"\">"
	}

	return result
}


func formWidgetGenerator(data interface{}, widgetSize int64, widgetTitle string, datatable string) template.HTML {
	form := "<div class=\"col-md-" + strconv.FormatInt(widgetSize, 10) + "\"><div class=\"panel panel-default\">" +
		"<div class=\"panel-heading\"><h3 class=\"panel-title\">" + widgetTitle + "</h3></div>" +
		"<div class=\"panel-body\">"

	form += formGenerator(data, datatable)

	form += "</div></div></div>"
	return template.HTML([]byte(form))
}

func tableWithFormWG(data interface{}, widgetSize int64, widgetTitle string, datatable string) template.HTML {
	mwin := "<div id=\"modal" + strings.ToLower(strings.Replace(widgetTitle, " ", "", -1)) + "\" class=\"modal fade\" role=\"dialog\">" +
		"<div class=\"modal-dialog\"><div class=\"modal-content\"><div class=\"modal-header\">" +
		"<button type=\"button\" class=\"close\" data-dismiss=\"modal\">&times;</button>" +
		"<h4 class=\"modal-title\"> Add " + widgetTitle + "</h4></div><div class=\"modal-body\">"

	mwin += formGenerator(data, datatable)

	mwin += "</div></div></div></div>"

	twin := "<div class=\"col-md-" + strconv.FormatInt(widgetSize, 10) + "\"><div class=\"panel panel-default\"><div class=\"panel-heading\"><h3 class=\"panel-title\">" + widgetTitle + "</h3>" +
		"<button type=\"button\" class=\"btn btn-default waves-effect m-b-6\" data-toggle=\"modal\" data-target=\"#modal" + strings.ToLower(strings.Replace(widgetTitle, " ", "", -1)) + "\">Add Data</button>" +
		"</div><div class=\"panel-body\"><div class=\"row\"><div class=\"col-md-12 col-sm-12 col-xs-12\">"
	twin += tableGenerator(data)
	twin += "</div></div></div></div></div>"
	return template.HTML([]byte(mwin + twin))
}

//======================================================================================================================
// Generate table
func tableGeneratorWith(data interface{}, datatable string) string {
	tb := "<table datasrc=\"" + datatable + "\" class=\"table table-bordered table-striped\" id=\"datatable-editable\"><thead><tr>"

	// add table header
	preF := reflect.TypeOf(data).Elem()
	for i := 0; i < preF.NumField(); i++ {
		if preF.Field(i).Name != "ID" {
			tb += "<th>" + preF.Field(i).Name + "</th>"
		}
	}
	tb += "<th></th></tr></thead><tbody>"

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
			if r, ok := pre.Field(i).Interface().(bool); ok {
				value = strconv.FormatBool(r)
			}
			if r, ok := pre.Field(i).Interface().(bson.ObjectId); ok {
				if preF.Field(i).Name == "Groupid" {
					for _, rowS := range dataSource.LoadDeviceGroup() {
						if rowS.ID == r {
							value = rowS.Name
						}
					}
				}
				if preF.Field(i).Name == "ID" {
					value = r.Hex()
				}
			}
			if i == 0 {
				tmp = "<tr class=\"gradeX\" id=\"" + value + "\">"
			} else {
				tmp = tmp + "<td>" + value + "</td>"
			}

		}
		tb = tb + tmp + "<td class=\"actions\"><a href=\"#\" class=\"hidden on-editing save-row\"><i class=\"fa fa-save\"></i></a>" +
			"<a href=\"#\" class=\"hidden on-editing cancel-row\"><i class=\"fa fa-times\"></i></a>" +
			"<a href=\"#\" class=\"on-default edit-row\"><i class=\"fa fa-pencil\"></i></a>" +
			"<a href=\"#\" class=\"on-default remove-row\"><i class=\"fa fa-trash-o\"></i></a></td></tr>"
	}
	tb += "</tbody></table>"
	return tb
}

//Generate Table Widget Editale Table
func editableTableWidgetGenerate(data interface{}, widgetSize int64, widgetTitle string, datatable string) template.HTML {
	editTableWg := "<div class=\"col-sm-" + strconv.FormatInt(widgetSize, 10) + "\"><div class=\"row\"><h4 class=\"pull-left page-title\">" + widgetTitle + "</h4></div>" +
		"<div class=\"panel\"><div class=\"panel-body\"> <div class=\"row\"><div class=\"col-sm-6\"><div class=\"m-b-30\">" +
		"<button id=\"addToTable\" class=\"btn btn-primary waves-effect waves-light\">Add <i class=\"fa fa-plus\"></i></button>" +
		"</div></div></div>" + tableGeneratorWith(data, datatable) + "</div></div></div>" +
		"<div id=\"dialog\" class=\"modal-block mfp-hide\"><section class=\"panel panel-info panel-color\">" +
		"<header class=\"panel-heading\"><h2 class=\"panel-title\">Are you sure?</h2></header>" +
		"<div class=\"panel-body\"><div class=\"modal-wrapper\"><div class=\"modal-text\">" +
		"<p>Are you sure that you want to delete this row?</p></div></div><div class=\"row m-t-20\">" +
		"<div class=\"col-md-12 text-right\"><button id=\"dialogConfirm\" class=\"btn btn-primary\">Confirm</button>" +
		"<button id=\"dialogCancel\" class=\"btn btn-default\">Cancel</button></div></div></div></section></div>"

	return template.HTML([]byte(editTableWg))
}

//======================================================================================================================
//Left Side bar Menu Generator
func MenuGenerator(data interface{}) template.HTML {
	menu := "<div id=\"sidebar-menu\"><ul>"

	dataSlice := reflect.ValueOf(data)
	if dataSlice.Kind() != reflect.Slice {
		panic("Data is not a slice")
	}
	vdata := make([]interface{}, dataSlice.Len())
	for i := 0; i < dataSlice.Len(); i++ {
		vdata[i] = dataSlice.Index(i).Interface()
	}
	childMenu := dataLoader.ChildMenuList()

	for _, v := range vdata {

		pre := reflect.ValueOf(v)
		if r, ok := pre.Interface().(datasource.MenuGroups); ok {
			//
			childs := checkContinueValueChild(childMenu, r.ID)
			if len(childs) == 0 {
				menu += "<li><a href=\"#\" class=\"waves-effect\"><i class=\"md md-home\"></i><span> " + r.Title + " </span></a></li>"
			} else {
				menu += "<li class=\"has_sub\"> <a href=\"#\" class=\"waves-effect\">" +
					"<i class=\"md md-mail\"></i><span> " + r.Title + " </span><span class=\"pull-right\"><i class=\"md md-add\"></i></span></a>" +
					"<ul class=\"list-unstyled\">"
				for _, subMenu := range childs {
					menu += "<li><a href=\"/page?id=" + subMenu.MonitoringPagesID + "\">" + subMenu.Title + "</a></li>"
				}
				menu += "</ul></li>"
			}
		}
	}
	menu += "<li><a href=\"settings\" class=\"waves-effect\"><i class=\"md  md-settings\"></i><span> Settings </span></a></li>"
	menu += "</ul><div class=\"clearfix\"></div></div>"
	return template.HTML([]byte(menu))
}

func checkContinueValueChild(v1 []datasource.ChildMenu, value bson.ObjectId) []datasource.ChildMenu {
	var childObjects []datasource.ChildMenu
	for _, v := range v1 {
		if v.MenuGroupID == value.Hex() {
			childObjects = append(childObjects, v)
		}
	}
	return childObjects
}

func (self *WidgetListCreat) WidgetGenerate(data interface{}, widgetSize int64, widgetTitle string, widgetType string, datatable string) Widget {
	dataSource = datasource.MonitoringBase{}

	var wg Widget

	switch widgetType {
	case "table":
		wg = &ReadyWidget{tableWidgetGenerator(data, widgetSize, widgetTitle)}
	case "tablein":
		wg = &ReadyWidget{tableWithFormWG(data, widgetSize, widgetTitle, datatable)}
	case "form":
		wg = &ReadyWidget{formWidgetGenerator(data, widgetSize, widgetTitle, datatable)}
	case "etable":
		wg = &ReadyWidget{editableTableWidgetGenerate(data, widgetSize, widgetTitle, datatable)}
	default:
		log.Fatalln("Unknown Error")
	}

	self.registerWidget(wg)
	return wg
}
