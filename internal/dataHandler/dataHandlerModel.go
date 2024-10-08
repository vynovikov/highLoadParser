package dataHandler

type KeyDetailed struct {
	Ts   string
	Part int
}

type HeaderData struct {
	FormName    string
	FileName    string
	HeaderBytes []byte
}

type Value struct {
	H HeaderData `json:"h"`
	E int        `json:"e"`
}
