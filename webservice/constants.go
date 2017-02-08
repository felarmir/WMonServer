package webservice

// Widget Types
const (
	EditbleTable = "editable"
	TableWithForm = "formtable"
	BasicTable = "table"
	FormWidget = "formwidget"
)

// widget map 
func WidgetTypeMap() map[int]string {
	widgetsmap := map[int]string {
		0 : EditbleTable,
		1 : TableWithForm,
		2 : BasicTable,
		3 : FormWidget,
	}
	return widgetsmap
}
