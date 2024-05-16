package dataHandler

type DataHandler interface {
	Create(DataHandlerDTO) error
	Read(DataHandlerDTO) (value, error)
	Updade(DataHandlerDTO) error
	Delete(string) error
	//	Check(DataHandlerDTO) (Presence, error)
}
