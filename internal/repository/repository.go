package repository

import (
	"github.com/vynovikov/highLoadParser/internal/dataHandler"
)

type ParserRepository interface {
	Check(dataHandler.DataHandlerDTO) (dataHandler.Presence, error)
}

type repositoryStruct struct {
	dataHandler dataHandler.DataHandler
}

func NewParserRepository(dh dataHandler.DataHandler) *repositoryStruct {
	return &repositoryStruct{
		dataHandler: dh,
	}
}

func (r *repositoryStruct) Check(d dataHandler.DataHandlerDTO) (dataHandler.Presence, error) {

	return r.dataHandler.Check(d)
}
