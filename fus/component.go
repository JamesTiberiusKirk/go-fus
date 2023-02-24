package fus

import (
	"log"

	"github.com/JamesTiberiusKirk/go-fus/renderer"
	"github.com/labstack/echo/v4"
)

type ComponentInterface interface {
	GetID() string
	GetTemplate() string
	GenerateComponentData(parentData interface{}) (echo.Map, error)
	GetCompoents() map[string]ComponentInterface
}

type Component struct {
	context         echo.Context
	id              string
	template        string
	getCompoentData func() (interface{}, error)
	compoents       map[string]ComponentInterface
}

type dataGetterFunc func() (interface{}, error)

func NewComponent(c echo.Context, id, template string,
	dataGetter dataGetterFunc) *Component {
	return &Component{
		context:         c,
		id:              id,
		template:        template,
		getCompoentData: dataGetter,
	}
}

func (comp *Component) GetTemplate() string {
	return comp.template
}

func (comp *Component) GetID() string {
	return comp.id
}

func (comp *Component) GenerateComponentData(parentData interface{}) (echo.Map, error) {
	log.Println("Setting frame as include")
	log.Println(comp.template)
	comp.context.Set(renderer.FrameEchoContextName, renderer.Include)

	echoData := echo.Map{
		// "meta": comp.buildComponentMetaData(comp.context),
		// "routes": routesMap,
	}

	cmpData, err := comp.getCompoentData()
	echoData["cmpData"] = cmpData
	echoData["error"] = err
	echoData["parentData"] = parentData

	// err = comp.context.Render(httpStatus, comp.template, echoData)
	// if err != nil {
	// 	return echoData, err
	// }

	return echoData, nil
}

func (comp *Component) GetCompoents() map[string]ComponentInterface {
	return comp.compoents
}

// func (comp *Component) buildComponentMetaData(c echo.Context) MetaData {
// 	return MetaData{
// 		ID: comp.id,
// 		// Title:    comp.Title,
// 		URLError: c.QueryParam("error"),
// 		Success:  c.QueryParam("success"),
// 	}
// }
