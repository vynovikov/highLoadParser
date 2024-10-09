package dataHandler

import "errors"

var (
	errHeaderEnding   error = errors.New("ending of header")
	errHeaderNotFull  error = errors.New("header is not full")
	errHeaderNotFound error = errors.New("header is not found")
	ErrKeyNotFound    error = errors.New("key is not found")
)

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

type Value1 struct {
	H HeaderData1 `json:"headerData"`
	E int         `json:"isEndingNeeded"`
}

type Value struct {
	H HeaderData `json:"headerData"`
	E int        `json:"isEndingNeeded"`
}
