package infrastructure

import "github.com/vynovikov/highLoadParser/internal/repository"

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

type DataPiece1 interface {
	Part() int
	TS() string
	Body() []byte
	B() repository.Disposition
	E() repository.Disposition
}

type InfrastructureUnit struct {
	part int
	ts   string
	body []byte
	b    Disposition
	e    Disposition
}

func NewInfrastructureUnit(d DataPiece) *InfrastructureUnit {
	return &InfrastructureUnit{
		part: d.Part(),
		ts:   d.TS(),
		body: d.Body(),
		b:    Disposition(d.B()),
		e:    Disposition(d.E()),
	}
}

func (i *InfrastructureUnit) Part() int {
	return i.part
}

func (i *InfrastructureUnit) TS() string {
	return i.ts
}
func (i *InfrastructureUnit) Body() []byte {
	return i.body
}
func (i *InfrastructureUnit) B() repository.Disposition {
	return repository.Disposition(i.b)
}
func (i *InfrastructureUnit) E() repository.Disposition {
	return repository.Disposition(i.e)
}

type TransferUnitStruct struct {
	TH TransferHeader
	TB TransferBody
}

type TransferHeader struct {
}

type TransferBody struct {
	B []byte
}

type TransferUnit interface {
	Tx() error
}

func (t *TransferUnitStruct) Tx() error {

	return nil
}

type Presence struct {
}
