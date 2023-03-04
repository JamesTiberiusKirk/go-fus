package fhtml

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestElem(t *testing.T) {
	t.Run("Successfully testing", func(t *testing.T) {
		t.Run("one element", func(t *testing.T) {
			actual := Elem(
				&Opts{Tag: "div", ID: "testID"},
				"test text in the div",
			)
			expected := "<div id='testID'>\n\ttest text in the div\n</div>"
			assert.Equal(t, expected, actual)
		})
		t.Run("multiple nested elements", func(t *testing.T) {
			actual := Elem(
				&Opts{
					Tag: "div",
					ID:  "testID",
					CSS: ElemCSS{
						Border: ToString("red, 2px, solid"),
					},
				},
				Elem(&Opts{Tag: "b"}, "Test text inside the element"),
			)

			//nolint:lll // Dont care about it in tests.
			expected := "<div id='testID' style='border:red, 2px, solid;'>\n\t<b>\n\t\tTest text inside the element\n\t</b>\n</div>"

			assert.Equal(t, expected, actual)
		})
	})
}
