package fhtml

import (
	"bytes"
	"reflect"
	"strings"
	"text/template"
)

const (
	//nolint:lll // tmpl string
	elemTemplate = "<{{.Tag}}{{ if .ID }} id='{{.ID}}'{{ end }}{{if .Classes}} class='{{.Classes}}'{{end}}{{.Attribs}}{{if .CSS}} style='{{.CSS}}'{{end}}>\n\t{{.Inner}}\n</{{.Tag}}>"

	cssTagName = "css"
)

type elemStructData struct {
	Tag     string
	ID      string
	Attribs string
	Inner   string
	CSS     string
	Classes string
}

type Opts struct {
	Tag          string
	ID           string
	Class        []string
	CSS          *ElemCSS
	OtherAttribs map[string]string
}

func Elem(attribs Opts, innerHTML ...string) string {
	var b bytes.Buffer

	inner := ""
	for _, i := range innerHTML {
		inner += i
	}

	data := elemStructData{
		Tag:   attribs.Tag,
		Inner: indentHTMLOneLevel(inner),
	}

	if attribs.ID != "" {
		data.ID = attribs.ID
	}

	if attribs.CSS != nil {
		data.CSS = getCSSStyle(*attribs.CSS)
	}

	for _, class := range attribs.Class {
		data.Classes += class + " "
	}

	tmpl, err := template.New("element").Parse(elemTemplate)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(&b, data)
	if err != nil {
		panic(err)
	}
	return b.String()
}

func ToString(str string) *string {
	return &str
}

func indentHTMLOneLevel(html string) string {
	html = strings.ReplaceAll(html, "\n", "\n\t")
	return html
}

func getCSSStyle(css ElemCSS) string {
	result := ""

	// Getting other attributes
	for k, v := range css.Other {
		result += k + ":" + v + ";"
	}

	val := reflect.ValueOf(css)

	// All of the fields apart from the ones without the css tag
	for i := 0; i < val.NumField(); i++ {
		fieldTag := val.Type().Field(i).Tag.Get(cssTagName)
		fieldValueInterface := val.Field(i).Interface()

		fieldValueString, ok := fieldValueInterface.(*string)
		if !ok || fieldValueString == nil || fieldTag == "" {
			continue
		}

		result += fieldTag + ": " + *fieldValueString + ";"
	}

	return result
}
