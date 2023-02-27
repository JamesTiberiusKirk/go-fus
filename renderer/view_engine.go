package renderer

import (
	"fmt"
	"html/template"
	"io"
	"sync"

	"github.com/labstack/echo/v4"
)

const (
	Include = "include"
)

const (
	FrameEchoContextName = "frame"
)

// ViewEngine view template engine.
type ViewEngine struct {
	config      Config
	tplMap      map[string]*template.Template
	tplMutex    sync.RWMutex
	fileHandler FileHandler
}

// New new template engine.
func New(config Config) *ViewEngine {
	if config.FileHandlerType == "" {
		config.FileHandlerType = SingleFolder
	}

	return &ViewEngine{
		config:      config,
		tplMap:      make(map[string]*template.Template),
		tplMutex:    sync.RWMutex{},
		fileHandler: getFileHandler(config.FileHandlerType),
	}
}

// Render render func called by echo c.Render().
func (e *ViewEngine) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	//nolint:errcheck // No error to check
	frame := c.Get(FrameEchoContextName).(string)
	if frame == "" {
		frame = Include
	}

	err := e.executeTemplate(w, name, data, frame)
	if e.config.Dev {
		return returnErrToBrowser(w, c, err)
	}

	return err
}

func (e *ViewEngine) executeTemplate(out io.Writer, name string, data interface{},
	frame string) error {
	var tpl *template.Template
	var err error

	allFuncs := make(template.FuncMap, 0)

	allFuncs["cmp"] = e.getCmpFunc(data)
	allFuncs["include"] = e.getIncludeFunc(data)
	allFuncs["includeJs"] = e.getIncludeJSFunc(data)
	allFuncs["includeTs"] = e.getIncludeTSFunc(data)

	// Get the plugin collection
	for k, v := range e.config.Funcs {
		allFuncs[k] = v
	}

	e.tplMutex.RLock()
	tpl, tplMapOk := e.tplMap[name]
	e.tplMutex.RUnlock()

	exeName := name

	framePath, ok := e.config.Frames[frame]
	if !ok && frame != Include {
		return fmt.Errorf("frame type not found %s", frame)
	}

	if framePath != "" {
		exeName = framePath
	}

	if !tplMapOk || e.config.DisableCache {
		tplList := make([]string, 0)

		tplList = appendIfNotEmpty(tplList, framePath)
		tplList = appendIfNotEmpty(tplList, name)
		// tplList = append(tplList, e.config.Partials...)

		// Loop through each template and test the full path
		tpl = template.New(name).Funcs(allFuncs).Delims(e.config.Delims.Left, e.config.Delims.Right)
		for _, t := range tplList {
			var data string
			data, err = e.fileHandler(e.config, t)
			if err != nil {
				return fmt.Errorf("error getting file template data %w", err)
			}
			var tmpl *template.Template
			if t == name {
				tmpl = tpl
			} else {
				tmpl = tpl.New(t)
			}
			_, err = tmpl.Parse(data)
			if err != nil {
				return fmt.Errorf("ViewEngine render parser name:%v, error: %w", t, err)
			}
		}

		e.tplMutex.Lock()
		e.tplMap[name] = tpl
		e.tplMutex.Unlock()
	}

	// Display the content to the screen
	err = tpl.Funcs(allFuncs).ExecuteTemplate(out, exeName, data)
	if err != nil {
		return fmt.Errorf("ViewEngine execute template error: %w", err)
	}

	return nil
}

func appendIfNotEmpty(array []string, value string) []string {
	if value != "" {
		array = append(array, value)
	}
	return array
}
