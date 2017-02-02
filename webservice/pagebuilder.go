package webservice

import (
	"html/template"
)

type PageDataBuilder interface {
	registerTableWidget(widget template.HTML)
	registerFormWidget(widget template.HTML)
}

type PageData struct {
	TableWidget []template.HTML
	FormWidget []template.HTML
	Menu template.HTML
	Tablescripts	bool
	ChartScripts	bool
}

func (self *PageData) registerTableWidget(widget template.HTML) {
	self.TableWidget = append(self.TableWidget, widget)
}

func (self *PageData) registerFormWidget(widget template.HTML) {
	self.FormWidget = append(self.FormWidget, widget)
}