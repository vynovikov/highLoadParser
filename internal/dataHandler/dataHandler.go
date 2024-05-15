package dataHandler

type DataHandler interface {
	Create(DataHandlerDTO) error
	Read(DataHandlerDTO) (value, error)
	Updade(DataHandlerDTO) error
	Delete(DataHandlerDTO) error
	Check(DataHandlerDTO) (Presence, error)
}
