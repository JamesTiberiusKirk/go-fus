package renderer

import (
	"html/template"
	io "io"

	"github.com/labstack/echo/v4"
)

const (
	devErrorTemplate = `
	<style>
	#error{
		border: solid red 5px;
		padding: 5px;
		margin: 5px;
	}
	</style>
	<div id="error">
		<h1 style="color:red">ERROR:</h1>
		<b>Request: </b> [{{.Method}}] {{.Path}}
		<br/>
		<b>Message:</b> {{.ErrorMessage}}
	</div>
	`
)

func returnErrToBrowser(w io.Writer, c echo.Context, returnError error) error {
	tmpl := template.New("error-return")
	_, err := tmpl.Parse(devErrorTemplate)
	if err != nil {
		return err
	}

	tmplData := struct {
		ErrorMessage string
		Method       string
		Path         string
	}{
		ErrorMessage: returnError.Error(),
		Method:       c.Request().Method,
		Path:         c.Request().URL.RequestURI(),
	}

	err = tmpl.Funcs(nil).ExecuteTemplate(w, "error-return", tmplData)
	if err != nil {
		return err
	}

	return nil
}
