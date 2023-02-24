package renderer

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/labstack/echo/v4"
)

type ComponentInterface interface {
	GetID() string
	GetTemplate() string
	GenerateComponentData(parentData interface{}) (echo.Map, error)
	GetCompoents() map[string]ComponentInterface
}

type (
	compnentFunc func(component ComponentInterface) (template.HTML, error)
	includeFunc  func(tmpl string) (template.HTML, error)
)

func (e *ViewEngine) getCmpFunc(data interface{}) compnentFunc {
	return func(component ComponentInterface) (template.HTML, error) {
		var html template.HTML

		componentData, errInclude := component.GenerateComponentData(data)
		if errInclude != nil {
			return html, fmt.Errorf("error generating component data for: %s, %w",
				component.GetID(), errInclude)
		}

		buf := new(bytes.Buffer)
		errInclude = e.executeTemplate(buf, component.GetTemplate(), componentData,
			Include)
		if errInclude != nil {
			return html, fmt.Errorf("error executing template for component: %s, %w",
				component.GetID(), errInclude)
		}

		//nolint:gosec // This is not user submitted data, we want that to be html because it is an include.
		html = template.HTML(buf.String())
		return html, nil
	}
}

func (e *ViewEngine) getIncludeJSFunc(data interface{}) includeFunc {
	return func(tmpl string) (template.HTML, error) {
		buf := new(bytes.Buffer)
		errInclude := e.executeTemplate(buf, tmpl, data, Include)
		//nolint:gosec // This is not user submitted data, we want that to be html because it is an include.
		js := template.JS(buf.String())
		//nolint:gosec // This is not user submitted data, we want that to be html because it is an include.
		return template.HTML("\n<script>\n" + js + "\n</script>\n"), errInclude
	}
}

func (e *ViewEngine) getIncludeTSFunc(data interface{}) includeFunc {
	return func(tmpl string) (template.HTML, error) {
		buf := new(bytes.Buffer)
		path := "jsdist/" + strings.Replace(tmpl, ".ts", ".js", 1)
		errInclude := e.executeTemplate(buf, path, data, Include)
		//nolint:gosec // This is not user submitted data, we want that to be html because it is an include.
		js := template.JS(buf.String())
		//nolint:gosec // This is not user submitted data, we want that to be html because it is an include.
		return template.HTML("\n<script>\n" + js + "\n</script>\n"), errInclude
	}
}

func (e *ViewEngine) getIncludeFunc(data interface{}) includeFunc {
	return func(tmpl string) (template.HTML, error) {
		buf := new(bytes.Buffer)
		errInclude := e.executeTemplate(buf, tmpl, data, Include)
		//nolint:gosec // This is not user submitted data, we want that to be html because it is an include.
		return template.HTML(buf.String()), errInclude
	}
}
