package dataHandler

import "github.com/vynovikov/study/highLoadParser/internal/logger"

type memoryDataHandlerStruct struct {
	Map    map[string]Value
	Buffer []DataPiece
}

func NewMemoryDataHandler() *memoryDataHandlerStruct {
	return &memoryDataHandlerStruct{}
}

func (m *memoryDataHandlerStruct) Create(d DataPiece) error {
	logger.L.Printf("in dataHandler creating dataPiece = %v\n", d)
	return nil
}

func (m *memoryDataHandlerStruct) Read(DataPiece) (Value, error) {
	return Value{}, nil
}

func (m *memoryDataHandlerStruct) Updade(DataPiece) error {
	return nil
}

func (m *memoryDataHandlerStruct) Delete(DataPiece) error {
	return nil
}
