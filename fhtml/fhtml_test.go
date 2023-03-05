package fhtml

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestElem(t *testing.T) {
	t.Run("Successfully testing", func(t *testing.T) {
		t.Run("one element", func(t *testing.T) {
			actual := Elem(
				Opts{Tag: "div"},
				"test text in the div",
			)
			expected := "<div>\n\ttest text in the div\n</div>"
			assert.Equal(t, expected, actual)
		})
		t.Run("one element with classes", func(t *testing.T) {
			actual := Elem(
				Opts{
					Tag: "div",
					Class: []string{
						"test",
						"test1",
					},
				},
				"test text in the div",
			)
			expected := "<div class='test test1 '>\n\ttest text in the div\n</div>"
			assert.Equal(t, expected, actual)
		})
		t.Run("one element with other css value", func(t *testing.T) {
			actual := Elem(
				Opts{
					Tag: "div",
					CSS: &ElemCSS{
						Other: map[string]string{
							"some-other-css-prop": "red",
						},
					},
				},
				"test text in the div",
			)
			expected := "<div style='some-other-css-prop:red;'>\n\ttest text in the div\n</div>"
			assert.Equal(t, expected, actual)
		})
		t.Run("one element with id", func(t *testing.T) {
			actual := Elem(
				Opts{Tag: "div", ID: "testID"},
				"test text in the div",
			)
			expected := "<div id='testID'>\n\ttest text in the div\n</div>"
			assert.Equal(t, expected, actual)
		})
		t.Run("one element with CSS", func(t *testing.T) {
			actual := Elem(
				Opts{
					Tag: "div",
					CSS: &ElemCSS{
						Border: ToString("red 2px solid"),
					},
				},
				"test text in the div",
			)
			expected := "<div style='border: red 2px solid;'>\n\ttest text in the div\n</div>"
			assert.Equal(t, expected, actual)
		})
		t.Run("multiple nested elements", func(t *testing.T) {
			actual := Elem(
				Opts{
					Tag: "div",
					ID:  "testID",
					CSS: &ElemCSS{
						Border: ToString("red 2px solid"),
					},
				},
				Elem(Opts{Tag: "b"}, "Test text inside the element"),
			)

			//nolint:lll // Dont care about it in tests.
			expected := "<div id='testID' style='border: red 2px solid;'>\n\t<b>\n\t\tTest text inside the element\n\t</b>\n</div>"

			assert.Equal(t, expected, actual)
		})
		t.Run("multiple nested elements", func(t *testing.T) {
			actual := Elem(
				Opts{
					Tag: "div",
				},
				Elem(Opts{Tag: "b"}, "Test text inside the element"),
			)

			expected := "<div>\n\t<b>\n\t\tTest text inside the element\n\t</b>\n</div>"

			assert.Equal(t, expected, actual)
		})
		t.Run("multiple repeated nested elements", func(t *testing.T) {
			actual := Elem(
				Opts{
					Tag: "div",
				},
				Elem(Opts{Tag: "b"}, "Test text inside the element"),
				Elem(Opts{Tag: "b"}, "Test text inside the element"),
				Elem(Opts{Tag: "b"}, "Test text inside the element"),
				Elem(Opts{Tag: "b"}, "Test text inside the element"),
			)

			//nolint:lll // Long lines dont care
			expected := "<div>\n\t<b>\n\t\tTest text inside the element\n\t</b><b>\n\t\tTest text inside the element\n\t</b><b>\n\t\tTest text inside the element\n\t</b><b>\n\t\tTest text inside the element\n\t</b>\n</div>"

			assert.Equal(t, expected, actual)
		})
		t.Run("multiple nested elements with CSS styles and ID", func(t *testing.T) {
			actual := Elem(
				Opts{
					Tag: "div",
					ID:  "testID",
					CSS: &ElemCSS{
						Border: ToString("red 2px solid"),
					},
				},
				Elem(Opts{Tag: "b"}, "Test text inside the element"),
			)

			//nolint:lll // Dont care about it in tests.
			expected := "<div id='testID' style='border: red 2px solid;'>\n\t<b>\n\t\tTest text inside the element\n\t</b>\n</div>"

			assert.Equal(t, expected, actual)
		})
	})
}
