package dataHandler

type DataHandler interface {
	Create(DataHandlerDTO, Boundary) (ProducerUnit, error)
	Read(DataHandlerDTO) (value, error)
	Updade(DataHandlerDTO, Boundary) (ProducerUnit, error)
	Delete(string) error
	Set(keyDetailed, value) error
}
