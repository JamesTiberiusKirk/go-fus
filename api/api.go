package api

type APIInterface interface {
}

type API struct {
	id     string
	routes RoutesInterface
}
