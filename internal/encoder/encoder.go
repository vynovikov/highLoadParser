package encoder

type Encoder interface {
	EncodeKey(TransferUnit) []byte
	EncodeValue(TransferUnit) []byte
}
