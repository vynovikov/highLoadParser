package service

import (
	"bytes"
	"encoding/json"

	"github.com/vynovikov/highLoadParser/internal/dataHandler"
	"github.com/vynovikov/highLoadParser/internal/entities"
	"github.com/vynovikov/highLoadParser/internal/infrastructure"
	"github.com/vynovikov/highLoadParser/internal/logger"
	"github.com/vynovikov/highLoadParser/pkg/byteOps"
)

type ParcerService interface {
	Serve(ParserServiceDTO)
}

type parcerServiceStruct struct {
	infrastructure infrastructure.Infrastructure
}

func NewParserService(i infrastructure.Infrastructure) *parcerServiceStruct {
	return &parcerServiceStruct{
		infrastructure: i,
	}
}

func (s *parcerServiceStruct) Serve(sDTO ParserServiceDTO) {

	logger.L.Infoln("in Serve got some data")

	sDTO.Evolve(0)

	bou := newDataHandlerBoundary(sDTO.Bou)

	for _, v := range sDTO.psus {

		serviceUnit := newServiceUnit(v)

		dhu := newDataHandlerUnit(serviceUnit)

		resTT, err := s.infrastructure.Register(dhu, bou)

		if err != nil {

			logger.L.Warn(err)
		}

		tsus := append(make([]infrastructure.TransferUnit, 0), newTransferUnit(resTT)...)

		s.infrastructure.Send(tsus)

	}

	if sDTO.pssu != nil {

		serviceUnit := newServiceUnit(sDTO.pssu)

		dhu := newDataHandlerUnit(serviceUnit)

		_, err := s.infrastructure.Register(dhu, bou)

		if err != nil {

			logger.L.Warn(err)
		}
	}

}

func newDataHandlerBoundary(boundary entities.Boundary) dataHandler.Boundary {

	return dataHandler.Boundary{
		Prefix: boundary.Prefix,
		Root:   boundary.Root,
		Suffix: boundary.Suffix,
	}
}

func (s *ParserServiceDTO) Evolve(start int) {

	if s == nil {

		return
	}

	b := s.Body[start:]

	boundaryCore := entities.GetBoundary(s.Bou)[2:]

	ll := make([]byte, 0, entities.MaxLineLimit)

	if len(b) < len(boundaryCore)+2*len(entities.Sep) {

		s.last = true

		psu := NewParserServiceUnit(NewParserServiceHeader(s.TS, s.Part, 1, 0), NewParserServiceBody(b))

		s.psus = append(s.psus, &psu)

		return
	}

	if start == 0 {

		be := make([]byte, 0)

		if len(b) > entities.MaxHeaderLimit {

			be = b[len(b)-entities.MaxLineLimit:]
		} else {

			be = b
		}

		lenbe := len(be)
		ll = entities.GetLineWithCRLFLeft(be, lenbe-1, entities.MaxLineLimit, s.Bou)

		if (len(ll) > 2 && byteOps.BeginningEqual(ll[2:], boundaryCore)) ||
			(len(ll) <= 2 && bytes.Contains([]byte(entities.Sep), ll)) { // last line equal to boundary begginning or vice versa

			if entities.IsLastBoundary(ll, []byte(""), s.Bou) { // last boundary in last line

				s.last = true
			} else {

				as := NewParserServiceSub(NewParserServiceSubHeader(s.TS, s.Part), NewParserServiceSubBody(ll))
				s.pssu = &as
			}

			b = b[:len(b)-len(ll)]

			s.Body = s.Body[:len(s.Body)-len(ll)]

			if bytes.Contains(b, []byte(entities.BoundaryField)) {

				bodyIdx := bytes.Index(b, boundaryCore)

				b = b[bodyIdx+len(entities.Sep)+len(boundaryCore):]

				start += bodyIdx + len(entities.Sep) + len(boundaryCore)
			}
		}
	}
	switch bytes.Count(b, boundaryCore) {
	case 0:
		psu := ParserServiceUnit{}

		if (s.pssu == nil && len(ll) > 2 && !byteOps.BeginningEqual(ll[2:], boundaryCore)) ||
			(s.pssu == nil && len(ll) <= 2 && !bytes.Contains([]byte(entities.Sep), ll)) ||
			(s.pssu == nil && s.last) {

			if s.last && start == 0 {

				psu = NewParserServiceUnit(NewParserServiceHeader(s.TS, s.Part, 1, 0), NewParserServiceBody(b))
			} else if s.last && start != 0 {

				psu = NewParserServiceUnit(NewParserServiceHeader(s.TS, s.Part, 0, 0), NewParserServiceBody(b))
			} else if !s.last && start == 0 {

				psu = NewParserServiceUnit(NewParserServiceHeader(s.TS, s.Part, 1, 1), NewParserServiceBody(b))
			} else if !s.last && start != 0 {

				psu = NewParserServiceUnit(NewParserServiceHeader(s.TS, s.Part, 0, 1), NewParserServiceBody(b))
			}

		} else if (s.pssu == nil && len(ll) > 2 && byteOps.BeginningEqual(ll[2:], boundaryCore)) ||
			(s.pssu == nil && len(ll) <= 2 && bytes.Contains([]byte(entities.Sep), ll)) {

			if start == 0 {

				psu = NewParserServiceUnit(NewParserServiceHeader(s.TS, s.Part, 1, 0), NewParserServiceBody(b))
			} else {

				psu = NewParserServiceUnit(NewParserServiceHeader(s.TS, s.Part, 0, 1), NewParserServiceBody(b))
			}

		} else {

			if start == 0 {

				psu = NewParserServiceUnit(NewParserServiceHeader(s.TS, s.Part, 1, 2), NewParserServiceBody(b))
			} else {

				psu = NewParserServiceUnit(NewParserServiceHeader(s.TS, s.Part, 0, 2), NewParserServiceBody(b))
			}
		}

		s.psus = append(s.psus, &psu)

	default:

		idx := bytes.Index(b, boundaryCore) + start

		psu := ParserServiceUnit{}

		if start == 0 {

			psu = NewParserServiceUnit(NewParserServiceHeader(s.TS, s.Part, 1, 0), NewParserServiceBody(b[:idx-start-len(entities.Sep)]))
		} else {

			psu = NewParserServiceUnit(NewParserServiceHeader(s.TS, s.Part, 0, 0), NewParserServiceBody(b[:idx-start-len(entities.Sep)]))
		}

		s.psus = append(s.psus, &psu)

		s.Evolve(idx + len(entities.Sep) + len(boundaryCore))

	}
}

func newTransferUnit(p dataHandler.ProducerUnit) []infrastructure.TransferUnit {

	res := make([]infrastructure.TransferUnit, 0, 2)

	headerMap := make(map[string]string)

	headerMap["formName"] = p.FormName()

	if len(p.FileName()) > 0 {

		headerMap["fileName"] = p.FileName()
	}

	val, err := json.Marshal(headerMap)
	if err != nil {

		logger.L.Warn(err)
	}
	if p.Start() {

		res = append(res, &infrastructure.TransferUnitStruct{
			I_key:   []byte("header"),
			I_value: val,
		})
	}

	res = append(res, &infrastructure.TransferUnitStruct{
		I_key:   []byte("body"),
		I_value: p.Body(),
	})

	if p.Final() { // TODO fix Final absense issue

		res = append(res, &infrastructure.TransferUnitStruct{
			I_key:   []byte("header"),
			I_value: []byte("final"),
		})
	}

	return res
}
