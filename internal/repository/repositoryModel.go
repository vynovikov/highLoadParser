package repository

import "github.com/vynovikov/highLoadParser/internal/dataHandler"

type Disposition int

const (
	False Disposition = iota
	True
	Probably
)

type RepositoryDTO interface {
	Part() int
	TS() string
	Body() []byte
	B() Disposition
	E() Disposition
}

type DataHandlerDTO interface {
	Part() int
	TS() string
	Body() []byte
	B() dataHandler.Disposition
	E() dataHandler.Disposition
}

type RepositoryDTOUnit struct {
	part int
	ts   string
	body []byte
	b    dataHandler.Disposition
	e    dataHandler.Disposition
}

func NewRepositoryDTOUnit(d RepositoryDTO) *RepositoryDTOUnit {
	return &RepositoryDTOUnit{
		part: d.Part(),
		ts:   d.TS(),
		body: d.Body(),
		b:    dataHandler.Disposition(d.B()),
		e:    dataHandler.Disposition(d.E()),
	}
}

func (d *RepositoryDTOUnit) Part() int {
	return d.part
}

func (d *RepositoryDTOUnit) TS() string {
	return d.ts
}

func (d *RepositoryDTOUnit) Body() []byte {
	return d.body
}

func (d *RepositoryDTOUnit) B() dataHandler.Disposition {
	return d.b
}

func (d *RepositoryDTOUnit) E() dataHandler.Disposition {
	return d.e
}

type TransferUnit struct {
}

type Presence struct {
}
