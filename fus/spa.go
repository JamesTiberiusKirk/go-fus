package fus

type SPA struct {
	MenuID  string
	Path    string
	Dist    string
	Index   string
	Routing bool
	Dev     SPADevParams
}

type SPADevParams struct {
	Host string
}
