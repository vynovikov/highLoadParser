package dataHandler

type DataHandler interface {
	Create(DataHandlerDTO, Boundary) error
	Read(DataHandlerDTO) (value, error)
	Updade(DataHandlerDTO) error
	Delete(string) error
	//	Check(DataHandlerDTO) (Presence, error)
}
