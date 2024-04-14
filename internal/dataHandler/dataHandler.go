package dataHandler

type DataHandler interface {
	Create(DataPiece) error
	Read(DataPiece) (Value, error)
	Updade(DataPiece) error
	Delete(DataPiece) error
}
