package renderer

import (
	"fmt"
	"html/template"
	"io"
	"log"
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

// M map interface for data.
type M map[string]interface{}

// New new template engine.
func New(config Config) *ViewEngine {
	return &ViewEngine{
		config:      config,
		tplMap:      make(map[string]*template.Template),
		tplMutex:    sync.RWMutex{},
		fileHandler: defaultFileHandler(),
	}
}

// Default new default template engine.
func Default() *ViewEngine {
	return New(DefaultConfig())
}

// Render render template for echo interface.
func (e *ViewEngine) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	//nolint:errcheck // No error to check
	frame := c.Get(FrameEchoContextName).(string)
	return e.RenderWriter(w, name, data, frame)
}

// RenderWriter render template with io.Writer.
func (e *ViewEngine) RenderWriter(w io.Writer, name string, data interface{},
	frame string) error {
	return e.executeTemplate(w, name, data, frame)
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
				log.Println("template equals to name", t)
				tmpl = tpl
			} else {
				log.Println("template not equals to name", t)
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

// FileHandler file handler interface.
type FileHandler func(config Config, tplFile string) (content string, err error)

func appendIfNotEmpty(array []string, value string) []string {
	if value != "" {
		array = append(array, value)
	}
	return array
}