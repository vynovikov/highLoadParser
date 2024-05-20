package dataHandler

type DataHandler interface {
	Create(DataHandlerDTO, Boundary) error
	Read(DataHandlerDTO) (value, error)
	Updade(DataHandlerDTO, Boundary) error
	Delete(string) error
}
