package entities

const (
	BoundaryField  = "boundary="
	Sep            = "\r\n"
	MaxLineLimit   = 100
	MaxHeaderLimit = 210
)

type Boundary struct {
	Prefix []byte
	Root   []byte
	Suffix []byte
}
