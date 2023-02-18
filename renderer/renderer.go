package renderer

import (
	"bytes"
	"fmt"
	"html/template"
	"io"

	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/labstack/echo/v4"
)

const (
	Master   TemplateType = "usesMaster"
	NoMaster TemplateType = "noMaster"
	Include  TemplateType = "include"
)

const (
	UseFrameName = "frame"
)

type TemplateType string

// DefaultConfig default config.
func DefaultConfig() Config {
	return Config{
		Root: "views",
		// Master:       "layouts/master",
		Partials:     []string{},
		Funcs:        make(template.FuncMap),
		DisableCache: false,
		Delims:       Delims{Left: "{{", Right: "}}"},
	}
}

// ViewEngine view template engine.
type ViewEngine struct {
	config      Config
	tplMap      map[string]*template.Template
	tplMutex    sync.RWMutex
	fileHandler FileHandler
}

// Config configuration options.
type Config struct {
	Root   string                  // view root
	Frames map[TemplateType]string // list of frames
	// Master       string            // template master
	// NoFrame      string            // template master
	Partials     []string         // template partial, such as head, foot
	Funcs        template.FuncMap // template functions
	DisableCache bool             // disable cache, debug mode
	Delims       Delims           // delimeters
}

// M map interface for data.
type M map[string]interface{}

// Delims delims for template.
type Delims struct {
	Left  string
	Right string
}

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
	frame := c.Get(UseFrameName).(bool)
	return e.RenderWriter(w, name, data, frame)
}

// RenderWriter render template with io.Writer.
func (e *ViewEngine) RenderWriter(w io.Writer, name string, data interface{}, useMaster bool) error {
	var frame TemplateType
	if useMaster {
		frame = Master
	} else {
		frame = NoMaster
	}

	return e.executeTemplate(w, name, data, frame)
}

//nolint:funlen,gocognit // Too much to unpack for this version
func (e *ViewEngine) executeTemplate(out io.Writer, name string, data interface{}, frame TemplateType) error {
	var tpl *template.Template
	var err error
	// var ok bool

	allFuncs := make(template.FuncMap, 0)

	allFuncs["include"] = func(tmpl string) (template.HTML, error) {
		buf := new(bytes.Buffer)
		errInclude := e.executeTemplate(buf, tmpl, data, Include)
		//nolint:gosec // This is not user submitted data, we want that to be html because it is an include.
		return template.HTML(buf.String()), errInclude
	}

	allFuncs["includeJs"] = func(tmpl string) (template.HTML, error) {
		buf := new(bytes.Buffer)
		errInclude := e.executeTemplate(buf, tmpl, data, Include)
		//nolint:gosec // This is not user submitted data, we want that to be html because it is an include.
		js := template.JS(buf.String())
		//nolint:gosec // This is not user submitted data, we want that to be html because it is an include.
		return template.HTML("\n<script>\n" + js + "\n</script>\n"), errInclude
	}

	allFuncs["includeTs"] = func(tmpl string) (template.HTML, error) {
		buf := new(bytes.Buffer)
		path := "jsdist/" + strings.Replace(tmpl, ".ts", ".js", 1)
		errInclude := e.executeTemplate(buf, path, data, Include)
		//nolint:gosec // This is not user submitted data, we want that to be html because it is an include.
		js := template.JS(buf.String())
		//nolint:gosec // This is not user submitted data, we want that to be html because it is an include.
		return template.HTML("\n<script>\n" + js + "\n</script>\n"), errInclude
	}

	// Get the plugin collection
	for k, v := range e.config.Funcs {
		allFuncs[k] = v
	}

	e.tplMutex.RLock()
	tpl, tplMapOk := e.tplMap[name]
	e.tplMutex.RUnlock()

	exeName := name

	framePath, ok := e.config.Frames[frame]
	if !ok {
		return fmt.Errorf("frame type not found %s", frame)
	}

	exeName = framePath

	if !tplMapOk || e.config.DisableCache {
		tplList := make([]string, 0)

		tplList = append(tplList, framePath)
		tplList = append(tplList, name)
		tplList = append(tplList, e.config.Partials...)

		// Loop through each template and test the full path
		tpl = template.New(name).Funcs(allFuncs).Delims(e.config.Delims.Left, e.config.Delims.Right)
		for _, t := range tplList {
			var data string
			data, err = e.fileHandler(e.config, t)
			if err != nil {
				return err
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

// FileHandler file handler interface.
type FileHandler func(config Config, tplFile string) (content string, err error)

// defaultFileHandler new default file handler
// This is just supposed to read the file and return a string as content.
func defaultFileHandler() FileHandler {
	return func(config Config, tplFile string) (string, error) {
		// Get the absolute path of the root template
		path, err := filepath.Abs(config.Root + string(os.PathSeparator) + tplFile)
		if err != nil {
			return "", fmt.Errorf("ViewEngine path:%v error: %w", path, err)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("ViewEngine render read name:%v, path:%v, error: %w", tplFile, path, err)
		}
		return string(data), nil
	}
}
