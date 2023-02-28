package fus

import (
	"bytes"
	"errors"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestViewEngine(t *testing.T) {
	ctrl := gomock.NewController(t)

	t.Run("Successful run", func(t *testing.T) {
		t.Run("with no test data and no frame name", func(t *testing.T) {
			ve := createTestViewEngine(t, nil, nil)

			mockContext := NewMockContext(ctrl)
			mockContext.EXPECT().
				Get(FrameEchoContextName).
				Return("")

			expectedHTML := "<html>\n  <head>\n  </head>\n  <body>\n  </body>\n</html>\n"

			var buf bytes.Buffer
			err := ve.Render(&buf, "full_page_no_data.gohtml", nil, mockContext)
			require.NoError(t, err)
			assert.Equal(t, expectedHTML, buf.String())
		})
		t.Run("with no test data", func(t *testing.T) {
			ve := createTestViewEngine(t, nil, nil)

			mockContext := NewMockContext(ctrl)
			mockContext.EXPECT().
				Get(FrameEchoContextName).
				Return("include")

			expectedHTML := "<html>\n  <head>\n  </head>\n  <body>\n  </body>\n</html>\n"

			var buf bytes.Buffer
			err := ve.Render(&buf, "full_page_no_data.gohtml", nil, mockContext)
			require.NoError(t, err)
			assert.Equal(t, expectedHTML, buf.String())
		})
		t.Run("with test data string", func(t *testing.T) {
			ve := createTestViewEngine(t, nil, nil)

			mockContext := NewMockContext(ctrl)
			mockContext.EXPECT().
				Get(FrameEchoContextName).
				Return("include")

			var buf bytes.Buffer

			testData := struct {
				Data string
			}{
				Data: "test",
			}

			expectedHTML := "<html>\n  <head>\n  </head>\n  <body>\n    test   \n  </body>\n</html>\n"

			err := ve.Render(&buf, "full_page_with_data.gohtml", testData, mockContext)
			require.NoError(t, err)
			assert.Equal(t, expectedHTML, buf.String())
		})

		t.Run("frame with empty page", func(t *testing.T) {
			frames := map[string]string{
				"normal-frame": "basic_frame_template.gohtml",
			}
			ve := createTestViewEngine(t, frames, nil)

			mockContext := NewMockContext(ctrl)
			mockContext.EXPECT().
				Get(FrameEchoContextName).
				Return("normal-frame")

			var buf bytes.Buffer

			expectedHTML := "<html>\n  <head>\n  </head>\n  <body>\n    \nPage\n   \n  </body>\n</html>\n"

			err := ve.Render(&buf, "page.gohtml", nil, mockContext)
			require.NoError(t, err)
			assert.Equal(t, expectedHTML, buf.String())
		})
		t.Run("frame with data on page", func(t *testing.T) {
			frames := map[string]string{
				"normal-frame": "basic_frame_template.gohtml",
			}
			ve := createTestViewEngine(t, frames, nil)

			mockContext := NewMockContext(ctrl)
			mockContext.EXPECT().
				Get(FrameEchoContextName).
				Return("normal-frame")

			var buf bytes.Buffer

			testData := struct {
				Data string
			}{
				Data: "test",
			}

			expectedHTML := "<html>\n  <head>\n  </head>\n  <body>\n    \nPage\n   \n  </body>\n</html>\n"

			err := ve.Render(&buf, "page.gohtml", testData, mockContext)
			require.NoError(t, err)
			assert.Equal(t, expectedHTML, buf.String())
		})
		t.Run("with custom funcmaps", func(t *testing.T) {
			tmplFunc := func(params string) template.HTML {
				return template.HTML(params + " lookin for this string")
			}

			frames := map[string]string{
				"normal-frame": "basic_frame_template.gohtml",
			}

			ve := createTestViewEngine(t, frames, template.FuncMap{
				"test": tmplFunc,
			})

			mockContext := NewMockContext(ctrl)
			mockContext.EXPECT().
				Get(FrameEchoContextName).
				Return("normal-frame")

			//nolint:lll // Its just an expected.
			expectedHTML := "<html>\n  <head>\n  </head>\n  <body>\n    \nPage\ntest lookin for this string\n   \n  </body>\n</html>\n"

			var buf bytes.Buffer
			err := ve.Render(&buf, "page_with_test_tmpl_func.gohtml", nil, mockContext)
			require.NoError(t, err)
			assert.Equal(t, expectedHTML, buf.String())
		})
		t.Run("in dev mode to return error to user", func(t *testing.T) {
			frameName := "normal-frame"
			frames := map[string]string{
				frameName: "basic_frame_template.gohtml",
			}

			ve := &viewEngine{
				config: viewEngineConfig{
					Root:   "testdata",
					Frames: frames,
					Delims: templateDelims{
						Left:  "{{",
						Right: "}}",
					},
					DisableCache: true,
					Dev:          true,
				},
				tplMap:   map[string]*template.Template{},
				tplMutex: sync.RWMutex{},
				fileHandler: func(config viewEngineConfig, tplFile string) (string, error) {
					path, err := filepath.Abs(config.Root + string(os.PathSeparator) + tplFile)
					require.NoError(t, err)

					data, err := os.ReadFile(path)
					require.NoError(t, err)

					return string(data), nil
				},
			}

			mockContext := NewMockContext(ctrl)
			mockContext.EXPECT().
				Get(FrameEchoContextName).
				Return("misspelled")
			mockContext.EXPECT().
				Request().
				Return(&http.Request{
					Method: http.MethodGet,
					URL: &url.URL{
						Path: "/",
					},
				})

			//nolint:lll // Dont care about test.
			expectedHTML := "\n\t<style>\n\t#error{\n\t\tborder: solid red 5px;\n\t\tpadding: 5px;\n\t\tmargin: 5px;\n\t}\n\t</style>\n\t<div id=\"error\">\n\t\t<h1 style=\"color:red\">ERROR:</h1>\n\t\t<b>Request: </b> [GET] /\n\t\t<br/>\n\t\t<b>Message:</b> frame type not found misspelled\n\t</div>\n\t"

			var buf bytes.Buffer
			err := ve.Render(&buf, "page_with_parse_err.gohtml", "", mockContext)
			require.NoError(t, err)
			assert.Equal(t, expectedHTML, buf.String())
		})
	})
	t.Run("Should faild when", func(t *testing.T) {
		t.Run("non existent frame option is passed", func(t *testing.T) {
			frames := map[string]string{
				"normal-frame": "basic_frame_template.gohtml",
			}

			ve := createTestViewEngine(t, frames, nil)

			mockContext := NewMockContext(ctrl)
			mockContext.EXPECT().
				Get(FrameEchoContextName).
				Return("nonexistent-frame")

			var buf bytes.Buffer
			err := ve.Render(&buf, "page_with_test_tmpl_func.gohtml", nil, mockContext)
			require.Equal(t, "", buf.String())
			assert.ErrorContains(t, err, "frame type not found nonexistent-frame")
		})
		t.Run("file handler error is returned", func(t *testing.T) {
			frames := map[string]string{
				"normal-frame": "basic_frame_template.gohtml",
			}

			expectedError := errors.New("test error")
			ve := &viewEngine{
				config: viewEngineConfig{
					Root:   "testdata",
					Frames: frames,
					Delims: templateDelims{
						Left:  "{{",
						Right: "}}",
					},
					DisableCache: true,
				},
				tplMap:   map[string]*template.Template{},
				tplMutex: sync.RWMutex{},
				fileHandler: func(config viewEngineConfig, tplFile string) (string, error) {
					return "", expectedError
				},
			}

			mockContext := NewMockContext(ctrl)
			mockContext.EXPECT().
				Get(FrameEchoContextName).
				Return("normal-frame")

			var buf bytes.Buffer
			err := ve.Render(&buf, "page_with_test_tmpl_func.gohtml", nil, mockContext)
			require.Equal(t, "", buf.String())
			assert.ErrorIs(t, err, expectedError)
		})
		t.Run("template parse error is produced", func(t *testing.T) {
			frames := map[string]string{
				"normal-frame": "basic_frame_template.gohtml",
			}

			ve := createTestViewEngine(t, frames, nil)

			mockContext := NewMockContext(ctrl)
			mockContext.EXPECT().
				Get(FrameEchoContextName).
				Return("normal-frame")

			var buf bytes.Buffer
			err := ve.Render(&buf, "page_with_test_tmpl_func.gohtml", nil, mockContext)
			require.Equal(t, "", buf.String())
			//nolint:lll // Dont care about linting in tests.
			assert.ErrorContains(t, err, "ViewEngine render parser name:page_with_test_tmpl_func.gohtml, error: template: page_with_test_tmpl_func.gohtml:3: function \"test\" not defined")
		})
	})
}

func createTestViewEngine(t *testing.T, frames map[string]string, funcs template.FuncMap) *viewEngine {
	ve := &viewEngine{
		config: viewEngineConfig{
			Root:   "testdata",
			Frames: frames,
			Delims: templateDelims{
				Left:  "{{",
				Right: "}}",
			},
			Funcs:        funcs,
			DisableCache: true,
		},
		tplMap:   map[string]*template.Template{},
		tplMutex: sync.RWMutex{},
		fileHandler: func(config viewEngineConfig, tplFile string) (string, error) {
			path, err := filepath.Abs(config.Root + string(os.PathSeparator) + tplFile)
			require.NoError(t, err)

			data, err := os.ReadFile(path)
			require.NoError(t, err)

			return string(data), nil
		},
	}
	return ve
}
