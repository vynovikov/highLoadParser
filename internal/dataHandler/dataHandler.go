package dataHandler

type DataHandler interface {
	Create(DataHandlerDTO, Boundary) (ProducerUnit, error)
	Read(DataHandlerDTO) (value, error)
	Updade(DataHandlerDTO, Boundary) (ProducerUnit, error)
	Delete(string) error
	Set(KeyDetailed, Value) error
	Get(KeyDetailed) (Value, error)
}
