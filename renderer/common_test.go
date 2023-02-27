package renderer

import (
	"html/template"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func createTestViewEngine(t *testing.T, frames map[string]string, funcs template.FuncMap) *ViewEngine {
	ve := &ViewEngine{
		config: Config{
			Root:   "testdata",
			Frames: frames,
			Delims: Delims{
				Left:  "{{",
				Right: "}}",
			},
			Funcs:        funcs,
			DisableCache: true,
		},
		tplMap:   map[string]*template.Template{},
		tplMutex: sync.RWMutex{},
		fileHandler: func(config Config, tplFile string) (string, error) {
			path, err := filepath.Abs(config.Root + string(os.PathSeparator) + tplFile)
			require.NoError(t, err)

			data, err := os.ReadFile(path)
			require.NoError(t, err)

			return string(data), nil
		},
	}
	return ve
}
