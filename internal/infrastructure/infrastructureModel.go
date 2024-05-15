package infrastructure

import "github.com/vynovikov/highLoadParser/internal/repository"

type Disposition int

const (
	False Disposition = iota
	True
	Probably
)

type InfraStructureDTO interface {
	Part() int
	TS() string
	Body() []byte
	B() Disposition
	E() Disposition
}

type RepositoryDTO interface {
	Part() int
	TS() string
	Body() []byte
	B() repository.Disposition
	E() repository.Disposition
}

type RepositoryDTOUnit struct {
	part int
	ts   string
	body []byte
	b    repository.Disposition
	e    repository.Disposition
}

func NewRepositoryDTOUnit(d InfraStructureDTO) *RepositoryDTOUnit {
	return &RepositoryDTOUnit{
		part: d.Part(),
		ts:   d.TS(),
		body: d.Body(),
		b:    repository.Disposition(d.B()),
		e:    repository.Disposition(d.E()),
	}
}

func (i *RepositoryDTOUnit) Part() int {
	return i.part
}

func (i *RepositoryDTOUnit) TS() string {
	return i.ts
}
func (i *RepositoryDTOUnit) Body() []byte {
	return i.body
}
func (i *RepositoryDTOUnit) B() repository.Disposition {
	return i.b
}
func (i *RepositoryDTOUnit) E() repository.Disposition {
	return i.e
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
