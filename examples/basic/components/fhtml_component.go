package components

import (
	"html/template"

	"github.com/JamesTiberiusKirk/go-fus/fhtml"
	"github.com/JamesTiberiusKirk/go-fus/fus"
	"github.com/labstack/echo/v4"
)

type FHTMLComponent struct {
	*fus.Component
}

func NewFHTMLComponent() *FHTMLComponent {
	return &FHTMLComponent{
		Component: fus.NewStringComponent(
			"fhtmlComponent",
			func(c echo.Context) (template.HTML, error) {
				rows := []string{}

				for i := 0; i < 100; i++ {
					row := fhtml.Elem(
						fhtml.Opts{
							Tag: "tr",
						},
						fhtml.Elem(
							fhtml.Opts{
								Tag: "td",
							},
							"table cell data",
						),
					)

					rows = append(rows, row)
				}

				htmlString := fhtml.Elem(
					fhtml.Opts{
						Tag: "div",
					},
					fhtml.Elem(
						fhtml.Opts{
							Tag: "b",
						},
						"This is hello world from fhtml\n",
					),
					fhtml.Elem(
						fhtml.Opts{
							Tag: "table",
							ID:  "tbaleID",
							CSS: &fhtml.ElemCSS{
								Border: fhtml.ToString("5px solid red"),
							},
						},
						fhtml.Elem(
							fhtml.Opts{
								Tag: "thead",
							},
							fhtml.Elem(
								fhtml.Opts{
									Tag: "tr",
								},
								fhtml.Elem(fhtml.Opts{Tag: "th"}, "Table cell name"),
							),
						),
						fhtml.Elem(
							fhtml.Opts{
								Tag: "tbody",
							},
							rows...,
						),
					),
				)

				//nolint:gosec // We dont want to escape html.
				return template.HTML(htmlString), nil
			},
		),
	}
}
