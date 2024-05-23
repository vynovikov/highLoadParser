package infrastructure

type TransferUnit interface {
	Key() []byte
	Value() []byte
}
