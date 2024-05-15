package repository

type Disposition int

const (
	False Disposition = iota
	True
	Probably
)

type DataPiece interface {
	Part() int
	TS() string
	Body() []byte
	B() Disposition
	E() Disposition
}

type RepositoryUnit struct {
	part int
	ts   string
	body []byte
	b    Disposition
	e    Disposition
}

func NewRepositoryUnit(d DataPiece) *RepositoryUnit {
	return &RepositoryUnit{
		part: d.Part(),
		ts:   d.TS(),
		body: d.Body(),
		b:    Disposition(d.B()),
		e:    Disposition(d.E()),
	}
}

func (i *RepositoryUnit) Part() int {
	return i.part
}

func (i *RepositoryUnit) TS() string {
	return i.ts
}
func (i *RepositoryUnit) Body() []byte {
	return i.body
}
func (i *RepositoryUnit) B() Disposition {
	return i.b
}
func (i *RepositoryUnit) E() Disposition {
	return i.e
}

type TransferUnit struct {
}

type Presence struct {
}
