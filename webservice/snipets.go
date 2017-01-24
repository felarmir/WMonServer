package webservice

import (
	"strconv"
	"reflect"
	"html/template"
	"fmt"
)

type Table struct {
	Tabledata template.HTML
}

// colSize betwen 3 and 10
func getHeaderDIV() string {
	snp := "<div class=\"col-md-6\"> <div class=\"content-box-large\"> " +
		"<div class=\"panel-heading\"> <div class=\"panel-title\">Border Table</div> "+
		"<div class=\"panel-options\"> <a href=\"#\" data-rel=\"collapse\"><i class=\"glyphicon glyphicon-refresh\"></i></a> "+
		"<a href=\"#\" data-rel=\"reload\"><i class=\"glyphicon glyphicon-cog\"></i></a> "+
		"</div> </div> <div class=\"panel-body\"> <table class=\"table table-bordered\"> "
	return snp
}

// generate thead
func generateTapleHeader(data interface{}) string {
	topen := "<thead> <tr>"
	tclose := "</tr></thead>"
	var header string


	fmt.Println(data)
	preF := reflect.ValueOf(data).Elem()
	fields := preF.Type()
	for i := 0; i < preF.NumField(); i++ {
		header = header + "<th>" + fields.Field(i).Name + "</th>"
	}
	return topen + header + tclose
}

// generate table with data
func generateTableData(data ...interface{}) string {
	result := "<tbody>"

	for _, v := range data {
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
		result  = result + tmp + "</tr>"
	}

	return result + "</tbody>"
}

func  TableGenerator(data interface{}) Table {
	var result string
	result = getHeaderDIV() + generateTapleHeader(data)// + generateTableData(data)
	result += "</table> </div> </div> </div>"
	return Table{template.HTML([]byte(result))}
}
