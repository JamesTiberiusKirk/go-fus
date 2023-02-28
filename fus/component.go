package fus

import (
	"errors"

	"github.com/JamesTiberiusKirk/go-fus/fusint"
	"github.com/labstack/echo/v4"
)

type Component struct {
	context         echo.Context
	id              string // For use in the template
	template        string
	getCompoentData dataGetterFunc
	compoents       []fusint.ComponentInterface
}

type dataGetterFunc func(c echo.Context, params interface{}) (interface{}, error)

func NewComponent(id, template string, dataGetter dataGetterFunc, compoents ...fusint.ComponentInterface) *Component {
	return &Component{
		id:              id,
		template:        template,
		getCompoentData: dataGetter,
		compoents:       compoents,
	}
}

func (comp *Component) GetTemplate() string {
	return comp.template
}

func (comp *Component) GetID() string {
	return comp.id
}

func (comp *Component) GenerateComponentData(parentData interface{}, params ...interface{}) (echo.Map, error) {
	if comp.context == nil {
		return nil, errors.New("no context provided")
	}

	echoData := echo.Map{}

	cmpData, err := comp.getCompoentData(comp.context, params[0])
	echoData["cmpData"] = cmpData
	echoData["error"] = err
	echoData["parentData"] = parentData

	components := comp.GetCompoents()
	compoentsMap := echo.Map{}
	for _, component := range components {
		component.SetContext(comp.context)
		compoentsMap[component.GetID()] = component
	}

	echoData["c"] = compoentsMap

	return echoData, nil
}

func (comp *Component) GetCompoents() []fusint.ComponentInterface {
	return comp.compoents
}

func (comp *Component) SetContext(c echo.Context) {
	comp.context = c
}
