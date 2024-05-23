package dataHandler

type DataHandler interface {
	Create(DataHandlerDTO, Boundary) (*TT, error)
	Read(DataHandlerDTO) (value, error)
	Updade(DataHandlerDTO, Boundary) (*TT, error)
	Delete(string) error
}
