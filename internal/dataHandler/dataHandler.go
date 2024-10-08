package dataHandler

type DataHandler interface {
	Create(DataHandlerDTO, Boundary) (ProducerUnit, error)
	Read(DataHandlerDTO) (Value, error)
	Updade(DataHandlerDTO, Boundary) (ProducerUnit, error)
	Delete(KeyDetailed) error
	Set(KeyDetailed, Value) error
	Get(KeyDetailed) (Value, error)
}
