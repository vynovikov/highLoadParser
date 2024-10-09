package dataHandler

type KeyDetailed struct {
	Ts   string
	Part int
}

type HeaderData struct {
	FormName string
	FileName string
	Header   []byte
}

type HeaderData1 struct {
	FormName string
	FileName string
	Header   string
}

type Value struct {
	H HeaderData1 `json:"headerData"`
	E int         `json:"isEndingNeeded"`
}
