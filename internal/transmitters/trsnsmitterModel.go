package transmitters

type TransferUnit interface {
	Tx() error
}
