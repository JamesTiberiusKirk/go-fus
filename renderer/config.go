package renderer

import "html/template"

// Config configuration options.
type Config struct {
	Root            string            // view root
	Frames          map[string]string // list of frames
	Partials        []string          // template partial, such as head, foot
	Funcs           template.FuncMap  // template functions
	DisableCache    bool              // disable cache, debug mode
	Delims          Delims            // delimeters
	FileHandlerType FileHandlerType   // type of existing file handler
	Dev             bool
}

// DefaultConfig default config.
// TODO: actually implement this?
// func DefaultConfig() Config {
// 	return Config{
// 		Root:         "views",
// 		Partials:     []string{},
// 		Funcs:        make(template.FuncMap),
// 		DisableCache: false,
// 		Delims:       Delims{Left: "{{", Right: "}}"},
// 	}
// }

// Delims delims for template.
type Delims struct {
	Left  string
	Right string
}
