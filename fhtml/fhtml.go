package fhtml

import (
	"bytes"
	"reflect"
	"strings"
	"text/template"
)

const (
	//nolint:lll // tmpl string
	elemTemplate = "<{{.T}}{{ if .ID }} id='{{.ID}}'{{ end }}{{.Attribs}}{{if .CSS}} style='{{.CSS}}'{{end}}>\n\t{{.Inner}}\n</{{.T}}>"
)

type elemStructData struct {
	Tag     string
	ID      string
	Attribs string
	Inner   string
	CSS     string
}

type Opts struct {
	Tag          string
	ID           string
	Class        []string
	CSS          ElemCSS
	OtherAttribs map[string]string
}

func Elem(attribs *Opts, innerHTML ...string) string {
	var b bytes.Buffer

	inner := ""
	for _, i := range innerHTML {
		inner += i
	}

	data := elemStructData{
		Tag:   attribs.Tag,
		Inner: indentHTMLOneLevel(inner),
	}

	if attribs != nil {
		data.ID = attribs.ID
		data.CSS = getCssStyle(attribs.CSS)
	}

	tmpl, _ := template.New("element").Parse(elemTemplate)
	_ = tmpl.Execute(&b, data)
	return b.String()
}

func ToString(str string) *string {
	return &str
}

func indentHTMLOneLevel(html string) string {
	html = strings.ReplaceAll(html, "\n", "\n\t")
	return html
}

func getCssStyle(css ElemCSS) string {
	result := ""

	// Getting other
	for k, v := range css.Other {
		result += k + ":" + v + ";"
	}

	val := reflect.ValueOf(css)

	// All of the fields apart from Other
	for i := 0; i < val.NumField(); i++ {
		fieldName := val.Type().Field(i).Name
		if fieldName == "Other" {
			break
		}

		fieldValueInterface := val.Field(i).Interface()
		fieldValueString := fieldValueInterface.(*string)

		result += fieldName + ":" + *fieldValueString + ";"
	}

	return result
}
