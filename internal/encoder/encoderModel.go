package encoder

type TransferUnit interface {
	TS() string
	Part() int
	FormName() string
	FileName() string
	Body() []byte
	Start() bool
	IsSub() bool
	End() bool
	Final() bool
}
