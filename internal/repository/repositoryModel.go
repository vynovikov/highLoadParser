package repository

type Disposition int

const (
	False Disposition = iota
	True
	Probably
)

/*
	type DTO struct {
		TS    string
		Value []DTOValue
	}

func NewDTO(ts string, val DTOValue) DTO {

		return DTO{
			TS:    ts,
			Value: []DTOValue{val},
		}
	}

type AppDistributoryUnit struct {
}

	type DTOValue struct {
		Part         int
		CreatedBySub bool
		Dispossition DTODisposition
	}

	func NewDTOValue(p int, cbs bool, dispo DTODisposition) DTOValue {
		return DTOValue{
			Part:         p,
			CreatedBySub: cbs,
			Dispossition: dispo,
		}
	}

	type DTODisposition struct {
		FormName        string
		FileName        string
		DispositionBody []byte
	}

	func NewDTODisposition(fo, fi string, dBody []byte) DTODisposition {
		return DTODisposition{
			FormName:        fo,
			FileName:        fi,
			DispositionBody: dBody,
		}
	}
*/
type DataPiece interface {
	Part() int
	TS() string
	B() Disposition
	E() Disposition
	Body() []byte
	Header() string
}
