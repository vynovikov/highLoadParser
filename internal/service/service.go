package service

import (
	"bytes"

	"github.com/vynovikov/highLoadParser/internal/entities"
	"github.com/vynovikov/highLoadParser/internal/logger"
	"github.com/vynovikov/highLoadParser/internal/repository"
	"github.com/vynovikov/highLoadParser/pkg/byteOps"
)

type ParcerService interface {
	Serve(ParserServiceDTO)
}

type parcerServiceStruct struct {
	repo repository.ParserRepository
}

func NewParserService(r repository.ParserRepository) *parcerServiceStruct {
	return &parcerServiceStruct{
		repo: r,
	}
}

func (s *parcerServiceStruct) Serve(sDTO ParserServiceDTO) {

	logger.L.Infoln("in Serve got some data")

	sDTO.Evolve(0)

	for _, v := range sDTO.psus {

		logger.L.Infof("psu header %s, body: %s\n", v.Header(), string(v.Body()))
	}

	if sDTO.pssu != nil {

		logger.L.Infof("pssu: %v\n", string(sDTO.pssu.Body()))
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

		psu := NewParserServiceUnit(NewParserServiceHeader(s.TS, s.Part, True, False), NewParserServiceBody(b))

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

				psu = NewParserServiceUnit(NewParserServiceHeader(s.TS, s.Part, True, False), NewParserServiceBody(b))
			} else if s.last && start != 0 {

				psu = NewParserServiceUnit(NewParserServiceHeader(s.TS, s.Part, False, False), NewParserServiceBody(b))
			} else if !s.last && start == 0 {

				psu = NewParserServiceUnit(NewParserServiceHeader(s.TS, s.Part, True, True), NewParserServiceBody(b))
			} else if !s.last && start != 0 {

				psu = NewParserServiceUnit(NewParserServiceHeader(s.TS, s.Part, False, True), NewParserServiceBody(b))
			}

		} else if (s.pssu == nil && len(ll) > 2 && byteOps.BeginningEqual(ll[2:], boundaryCore)) ||
			(s.pssu == nil && len(ll) <= 2 && bytes.Contains([]byte(entities.Sep), ll)) {

			if start == 0 {

				psu = NewParserServiceUnit(NewParserServiceHeader(s.TS, s.Part, True, False), NewParserServiceBody(b))
			} else {

				psu = NewParserServiceUnit(NewParserServiceHeader(s.TS, s.Part, False, True), NewParserServiceBody(b))
			}

		} else {

			if start == 0 {

				psu = NewParserServiceUnit(NewParserServiceHeader(s.TS, s.Part, True, Probably), NewParserServiceBody(b))
			} else {

				psu = NewParserServiceUnit(NewParserServiceHeader(s.TS, s.Part, False, Probably), NewParserServiceBody(b))
			}
		}

		s.psus = append(s.psus, &psu)

	default:

		idx := bytes.Index(b, boundaryCore) + start

		psu := ParserServiceUnit{}

		if start == 0 {

			psu = NewParserServiceUnit(NewParserServiceHeader(s.TS, s.Part, True, False), NewParserServiceBody(b[:idx-start-len(entities.Sep)]))
		} else {

			psu = NewParserServiceUnit(NewParserServiceHeader(s.TS, s.Part, False, False), NewParserServiceBody(b[:idx-start-len(entities.Sep)]))
		}

		s.psus = append(s.psus, &psu)

		s.Evolve(idx + len(entities.Sep) + len(boundaryCore))

	}
}
