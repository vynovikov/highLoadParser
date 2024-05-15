package dataHandler

type DataHandler interface {
	Create(DataPiece) error
	Read(DataPiece) (value, error)
	Updade(DataPiece) error
	Delete(DataPiece) error
	Check(DataPiece) (Presence, error)
}
