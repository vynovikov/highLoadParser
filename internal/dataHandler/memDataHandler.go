package dataHandler

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"

	"github.com/vynovikov/highLoadParser/internal/logger"
	"github.com/vynovikov/highLoadParser/pkg/byteOps"
	regexpops "github.com/vynovikov/highLoadParser/pkg/regexpOps"
)

type memoryDataHandlerStruct struct {
	Map    map[keyGeneral]map[keyDetailed]map[bool]value // two keys are for easy search
	Buffer []DataHandlerDTO
}

func NewMemoryDataHandler() *memoryDataHandlerStruct {
	return &memoryDataHandlerStruct{
		Map:    make(map[keyGeneral]map[keyDetailed]map[bool]value),
		Buffer: make([]DataHandlerDTO, 0),
	}
}

func (m *memoryDataHandlerStruct) Create(d DataHandlerDTO, bou Boundary) (ProducerUnit, error) {

	var l2Key bool

	kgen, kdet, resTT := newKeyGeneralFromDTO(d), newKeyDetailed(d), newResult(d)

	val, err := newValue(d, bou)
	if err != nil &&
		!errors.Is(err, errHeaderNotFull) &&
		!errors.Is(err, errHeaderEnding) {

		return resTT, err
	}

	if len(m.Map[kgen]) == 0 {

		l1, l2 := make(map[keyDetailed]map[bool]value), make(map[bool]value)

		if !d.IsSub() {

			kdet.part++

		} else {

			l2Key = true
		}

		l2[l2Key] = val

		l1[kdet] = l2

		if !d.Last() {

			m.Map[kgen] = l1
		}

		resTT.updateProducerUnit(d, val, len(val.H.headerBytes))

		return resTT, nil
	}

	// Not empty m.Map

	switch d.IsSub() {

	case false:

		if l1, ok := m.Map[kgen]; ok {

			if l2, ok := l1[kdet]; ok {

				delete(l1, kdet)

				if l3, ok := l2[false]; ok {

					if l3.E == 1 && d.E() == 1 {

						kdet.part++

						delete(l2, false)

					}

					delete(l2, false)

					l2[false] = val

				}
				if _, ok := l2[true]; ok && d.E() == 2 {

					kdet.part++

					l2[false] = val

				}

				l1[kdet] = l2

				m.Map[kgen] = l1

				resTT.updateProducerUnit(d, val, len(val.H.headerBytes))

				return resTT, nil
			}

			kdet.part++

			l1, l2 := make(map[keyDetailed]map[bool]value), make(map[bool]value)

			l2[l2Key] = val

			l1[kdet] = l2

			m.Map[kgen] = l1

			resTT.updateProducerUnit(d, val, len(val.H.headerBytes))

			return resTT, nil
		}

	default:

		l2Key = true

		if l1, ok := m.Map[kgen]; ok {

			if l2, ok := l1[kdet]; ok {

				if l3, ok := l2[false]; ok {

					if l3.E == 2 {

						delete(l1, kdet)

						kdet.part++

						l2[true] = val

					} else {

						l2[true] = val
					}
				}

				l1[kdet] = l2

				m.Map[kgen] = l1

				return resTT, nil
			}

			l1, l2 := make(map[keyDetailed]map[bool]value), make(map[bool]value)

			l2[l2Key] = val

			l1[kdet] = l2

			m.Map[kgen] = l1

			return resTT, nil
		}
	}

	return resTT, nil
}

func (m *memoryDataHandlerStruct) Read(DataHandlerDTO) (value, error) {

	return value{}, nil
}

func (m *memoryDataHandlerStruct) Updade(d DataHandlerDTO, bou Boundary) (ProducerUnit, error) {

	var (
		err error
	)

	kgen, kdet, oldValueFalseUpated, resTT, headerEndingIndex := newKeyGeneralFromDTO(d), newKeyDetailed(d), value{}, newResult(d), -1

	body := d.Body()

	if l1, ok := m.Map[kgen]; ok {

		if l2, ok := l1[kdet]; ok {

			l1New, l2New := make(map[keyDetailed]map[bool]value), make(map[bool]value)

			oldValueFalse := l2[false]

			if len(oldValueFalse.H.formName) == 0 {

				oldValueFalseUpated, _, err = updateValue(oldValueFalse, d, bou)
				if err != nil {

					logger.L.Warn(err)
				}

				headerEndingIndex = len(oldValueFalseUpated.H.headerBytes) - len(oldValueFalse.H.headerBytes)

			} else {

				oldValueFalseUpated = oldValueFalse

				oldValueFalseUpated.E = d.E()
			}

			if len(l2) > 1 {

				oldValueTrue := l2[true]

				oldHeader := oldValueTrue.H.headerBytes

				dispositionIndex := bytes.Index(body, []byte(CONTENT_DISPOSITION))

				if dispositionIndex > 0 && byteOps.SameByteSlices(append(oldHeader, body[:dispositionIndex-2]...), genBoundary(bou)) {

					d.SetBody(body[dispositionIndex:])

					val, err := newValue(d, bou)
					if err != nil {

						logger.L.Infoln(err)
					}

					l2New[false] = val

					delete(m.Map[kgen], kdet)

					kdet.part++

					l1New[kdet] = l2New

					m.Map[kgen] = l1New

					headerEndingIndex = len(val.H.headerBytes)

					resTT.updateProducerUnit(d, val, headerEndingIndex)

					return resTT, nil

				} else if oldValueFalse.E != 2 {

					if (len(l2) > 1 && d.E() == 2) ||
						(len(l2) == 1 && d.E() == 1) {

						kdet.part++
					}

					delete(l2, false)
					delete(m.Map[kgen], kdet)

					l2[false] = oldValueFalseUpated

					l1New[kdet] = l2

					m.Map[kgen] = l1New

					resTT.updateProducerUnit(d, oldValueFalseUpated, headerEndingIndex)

					return resTT, nil

				} else {

					if (len(l2) > 1 && d.E() == 2) ||
						(len(l2) > 1 && l2[false].E == 2 && l2[true].E == 2 && d.E() != 2) ||
						(len(l2) == 1 && d.E() == 1) {

						kdet.part++
					}

					delete(l2, true)
					delete(l2, false)
					delete(m.Map[kgen], kdet)

					l2[false] = oldValueFalseUpated

					l1New[kdet] = l2

					m.Map[kgen] = l1New

					d.SetBody(append(oldValueTrue.H.headerBytes, d.Body()...))

					resTT.updateProducerUnit(d, oldValueFalseUpated, 0)

					return resTT, nil
				}
			}

			if ok && d.E() == 0 {

				delete(m.Map[kgen], kdet)

				resTT.updateProducerUnit(d, oldValueFalseUpated, headerEndingIndex)

				return resTT, nil
			}

			if (len(l2) > 1 && d.E() == 2) ||
				(len(l2) == 1 && d.E() == 1) {

				kdet.part++
			}

			delete(l2, false)

			delete(m.Map[kgen], kdet)

			oldValueFalse.E = d.E()

			l2[false] = oldValueFalseUpated

			l1New[kdet] = l2

			m.Map[kgen] = l1New

			resTT.updateProducerUnit(d, oldValueFalseUpated, headerEndingIndex)

			return resTT, nil
		}
	}

	m.Buffer = append(m.Buffer, d)

	return nil, nil
}

func (m *memoryDataHandlerStruct) Delete(ts string) error {

	delete(m.Map, newKeyGeneralFromTS(ts))

	if len(m.Map) == 0 {

		m.Map = make(map[keyGeneral]map[keyDetailed]map[bool]value)
	}

	return nil
}

func (m *memoryDataHandlerStruct) Set(KeyDetailed, Value) error {
	return nil
}

func (m *memoryDataHandlerStruct) Get(KeyDetailed) (Value, error) {
	return Value{}, nil
}

// updates value returns index of first body symbol
func updateValue(val value, d DataHandlerDTO, bou Boundary) (value, int, error) {

	if len(val.H.formName) != 0 {

		return val, -1, nil
	}

	body, resValue := make([]byte, 0, maxHeaderLimit), value{}

	resValue.E = d.E()

	length := len(d.Body())

	if length >= maxHeaderLimit {

		body = append(body, d.Body()[:maxHeaderLimit]...)
	} else {

		body = append(body, d.Body()...)
	}

	headerEnding, err := getHeaderLines(body, bou)

	if err == nil ||
		errors.Is(err, errHeaderEnding) {

		headerFull := append(val.H.headerBytes, headerEnding...)

		resValue.H.headerBytes = headerFull

		resValue.H.formName, resValue.H.fileName = getFoFi(headerFull)

		return resValue, len(headerEnding), nil

	} else {

		return value{}, -1, err
	}
}

func newKeyGeneralFromDTO(d DataHandlerDTO) keyGeneral {

	return keyGeneral{
		ts: d.TS(),
	}
}

func newKeyGeneralFromTS(ts string) keyGeneral {

	return keyGeneral{
		ts: ts,
	}
}

func newValue(d DataHandlerDTO, bou Boundary) (value, error) {

	headerB, body, _ := make([]byte, 0, maxHeaderLimit), d.Body(), 0

	lengh := len(body)

	if lengh > maxHeaderLimit {

		headerB = append(headerB, d.Body()[:maxHeaderLimit]...)

	} else {

		headerB = append(headerB, d.Body()...)
	}

	exactHeaderBytes, err := getHeaderLines(headerB, bou)
	if err != nil {

		if errors.Is(err, errHeaderNotFull) ||
			errors.Is(err, errHeaderEnding) {

			return value{
				E: d.E(),
				H: headerData{
					headerBytes: exactHeaderBytes,
				},
			}, err
		}

		return value{}, err
	}

	fo, fi := getFoFi(exactHeaderBytes)

	return value{
		E: d.E(),
		H: headerData{
			formName:    fo,
			fileName:    fi,
			headerBytes: exactHeaderBytes,
		},
	}, nil
}

func newResult(d DataHandlerDTO) *ProducerUnitStruct {

	if d.IsSub() {

		return &ProducerUnitStruct{
			Dh_TS:    d.TS(),
			Dh_Part:  d.Part(),
			Dh_Body:  make([]byte, 0),
			Dh_IsSub: d.IsSub(),
		}
	}
	return &ProducerUnitStruct{
		Dh_TS:    d.TS(),
		Dh_Part:  d.Part(),
		Dh_Body:  make([]byte, 0, len(d.Body())),
		Dh_IsSub: d.IsSub(),
		Dh_End:   d.E() == 0,
		Dh_Final: d.Last(),
	}
}

func (t *ProducerUnitStruct) updateProducerUnit(d DataHandlerDTO, v value, headerEndingIndex int) {

	t.Dh_FormName = v.H.formName

	t.Dh_FileName = v.H.fileName

	if d.B() == 1 && headerEndingIndex > 0 {

		t.Dh_Body = d.Body()[headerEndingIndex:]

		t.Dh_Start = true

	} else if d.B() == 1 {

		t.Dh_Body = d.Body()

	} else if d.B() == 0 && len(v.H.formName) == 0 {

		t.Dh_Body = make([]byte, 0)

	} else if d.B() == 0 && len(v.H.formName) > 0 {

		t.Dh_Body = d.Body()[headerEndingIndex:]

		t.Dh_Start = true

	}
}

// getHeaderLines returns header lines found in b
// Tested in dataHandler_test.go
func getHeaderLines(b []byte, bou Boundary) ([]byte, error) {

	resL := make([]byte, 0)

	boundaryCore := genBoundary(bou)[2:]

	if len(b) == 0 {

		return resL, fmt.Errorf("zero len byte slice passed")
	}

	if b[0] == 10 { // preceding LF

		switch bytes.Count(b, []byte("\r\n")) {

		case 0: //  LF + rand

			resL = append(resL, b[0])

			return resL, fmt.Errorf("\"%s\" is %w", resL, errHeaderEnding)

		case 1: // LF + CRLF + rand

			resL = append(resL, b[0])
			resL = append(resL, []byte("\r\n")...)

			return resL, fmt.Errorf("\"%s\" is %w", resL, errHeaderEnding)

		case 2: // LF + CT + 2*CRLF + rand || LF + CDSuff + 2*CRLF + rand

			l0 := b[1:bytes.Index(b, []byte("\r\n"))]
			resL = append(resL, b[0])
			resL = append(resL, l0...)
			resL = append(resL, []byte("\r\n\r\n")...)

			return resL, fmt.Errorf("\"%s\" is %w", resL, errHeaderEnding)

		default: //  LF + CDinsuf + CRLF + CT + 2*CRLF + rand

			l0 := b[1:bytes.Index(b, []byte("\r\n"))]
			l1 := b[bytes.Index(b, []byte("\r\n"))+2 : byteOps.RepeatedIntex(b, []byte("\r\n"), 2)]

			if sufficientType(l0) == insufficient {

				resL = append(b[:1], l0...)
				resL = append(resL, []byte("\r\n")...)

				if regexpops.IsCTFull(l1) {

					resL = append(resL, l1...)
					resL = append(resL, []byte("\r\n\r\n")...)

					return resL, fmt.Errorf("\"%s\" is %w", resL, errHeaderEnding)
				}
			}

			resL = append(resL, b[0])

			return resL, fmt.Errorf("\"%s\" is %w", resL, errHeaderEnding)
		}

	}
	if b[len(b)-1] == 13 { // succeeding CR

		switch bytes.Count(b, []byte("\r\n")) {

		case 0: //  CD full + CR

			if sufficientType(b[:len(b)-1]) != incomplete {

				resL = append(resL, b...)

				return resL, fmt.Errorf("\"%s\" %w", resL, errHeaderNotFull)
			}

		case 1: // CDsuf + CRLF + CR || CDinsuf + CRLF + CT + CR

			l0 := b[:bytes.Index(b, []byte("\r\n"))]

			if sufficientType(l0) == sufficient {

				resL = append(resL, b...)

				return resL, fmt.Errorf("\"%s\" %w", resL, errHeaderNotFull)
			}

			if sufficientType(l0) == insufficient {

				resL = append(l0, []byte("\r\n")...)

				l1 := b[bytes.Index(b, []byte("\r\n"))+2 : len(b)-1]

				if regexpops.IsCTFull(l1) {

					resL = append(resL, l1...)
					resL = append(resL, []byte("\r")...)

					return resL, fmt.Errorf("\"%s\" %w", resL, errHeaderNotFull)
				}
			}

		case 2: // CDinsuf + CRLF + CT + CRLF + CR

			l0 := b[:bytes.Index(b, []byte("\r\n"))]
			l1 := b[bytes.Index(b, []byte("\r\n"))+2 : byteOps.RepeatedIntex(b, []byte("\r\n"), 2)]

			if sufficientType(l0) == insufficient {

				resL = append(l0, []byte("\r\n")...)

				if regexpops.IsCTFull(l1) {

					resL = append(resL, l1...)
					resL = append(resL, []byte("\r\n\r")...)

					return resL, fmt.Errorf("\"%s\" %w", resL, errHeaderNotFull)
				}
			}

		default: // CDinsuf + CRLF + CT + 2*CRLF + rand + CR

			l0 := b[:bytes.Index(b, []byte("\r\n"))]
			l1 := b[bytes.Index(b, []byte("\r\n"))+2 : byteOps.RepeatedIntex(b, []byte("\r\n"), 2)]

			if sufficientType(l0) == insufficient {

				resL = append(l0, []byte("\r\n")...)

				if regexpops.IsCTFull(l1) {

					resL = append(resL, l1...)
					resL = append(resL, []byte("\r\n\r\n")...)

					return resL, nil
				}
			}

			return nil, errHeaderNotFound
		}
	}
	// no precending LF no succeding CR

	switch bytes.Count(b, []byte("\r\n")) {

	case 0: // CD ->

		if regexpops.IsCDRight(b) {

			return b, fmt.Errorf("\"%s\" %w", b, errHeaderNotFull)
		}
		if isLastBoundaryPart(b, bou) {

			return b, nil
		}

		return nil, errHeaderNotFound

	case 1: // CD full + CRLF || CD full + CRLF + CT -> || CRLF || <-LastBoundary + CRLF || CRLF + Boundary-> || <-Boundary + CRLF

		l0 := b[:bytes.Index(b, []byte("\r\n"))]
		l1 := b[bytes.Index(b, []byte("\r\n"))+2:]

		if len(l0) == 0 && len(l1) > 0 && byteOps.BeginningEqual(boundaryCore, l1) {

			resL = append(resL, []byte("\r\n")...)
			resL = append(resL, l1...)

			return resL, fmt.Errorf("\"%s\" %w", resL, errHeaderNotFull)
		}

		if len(l0) == 0 {

			resL = append(l0, []byte("\r\n")...)

			return resL, fmt.Errorf("\"%s\" is %w", resL, errHeaderEnding)
		}

		if sufficientType(l0) == sufficient {

			resL = append(l0, []byte("\r\n")...)

			return resL, fmt.Errorf("\"%s\" %w", resL, errHeaderNotFull)
		}
		if sufficientType(l0) == insufficient {

			resL = append(l0, []byte("\r\n")...)

			if regexpops.IsCTRight(l1) {

				resL = append(resL, l1...)

				return resL, fmt.Errorf("\"%s\" %w", resL, errHeaderNotFull)
			}

		}
		if len(b) == bytes.Index(b, []byte("\r\n"))+2 { //last Boundary

			resL = append(resL, b...)

			return resL, fmt.Errorf("\"%s\" is the last", resL)
		}

		return nil, errHeaderNotFound

	case 2: // CD full insufficient + CRLF + CT full + CRLF || CD full sufficient + 2 CRLF + rand || CT full + 2 CRLF + rand || <-CT + 2CRLF + rand || 2 CRLF + rand

		l0 := b[:bytes.Index(b, []byte("\r\n"))]
		l1 := b[bytes.Index(b, []byte("\r\n"))+2 : byteOps.RepeatedIntex(b, []byte("\r\n"), 2)]

		if len(l0) == 0 { // on ending part is impossible, on beginning part: 2 * CRLF + rand || CRLF + rand + CRLF + rand

			resL = append(resL, []byte("\r\n")...)

			if len(l1) == 0 {

				resL = append(resL, []byte("\r\n")...)

				return resL, fmt.Errorf("\"%s\" is %w", resL, errHeaderEnding)
			}
			return resL, fmt.Errorf("\"%s\" is %w", resL, errHeaderEnding)
		}

		if sufficientType(l0) == sufficient { // on ending part CDSuf + 2 * CRLF || CDSuf + 2 * CRLF + rand, on beginning part CDSuf + 2 * CRLF + rand

			resL = append(l0, []byte("\r\n\r\n")...)

			return resL, nil
		}

		if sufficientType(l0) == insufficient { // on ending part CDInsuf + CRLF + CT + CRLF, on beginning part is impossible

			resL = append(l0, []byte("\r\n")...)

			if regexpops.IsCTFull(l1) {

				resL = append(resL, l1...)
				resL = append(resL, []byte("\r\n")...)

				return resL, fmt.Errorf("\"%s\" %w", resL, errHeaderNotFull)
			}
		}
		if regexpops.IsCDLeft(l0) && len(l1) == 0 { // on ending part is impossible, on beginning part <-CDsufficient + 2 * CRLF + rand

			resL = append(l0, []byte("\r\n\r\n")...)

			return resL, fmt.Errorf("\"%s\" is %w", resL, errHeaderEnding)
		}

		if regexpops.IsCTLeft(l0) && len(l1) == 0 { // on ending part is impossible, on beginning part <-CT + 2 * CRLF + rand

			resL = append(l0, []byte("\r\n\r\n")...)

			return resL, fmt.Errorf("\"%s\" is %w", resL, errHeaderEnding)
		}

		return nil, errHeaderNotFound

	default: // CD full insufficient + CRLF + CT full + 2*CRLF || CD full sufficient + 2*CRLF + rand + CRLF

		l0 := b[:bytes.Index(b, []byte("\r\n"))]
		l1 := b[byteOps.RepeatedIntex(b, []byte("\r\n"), 1)+2 : byteOps.RepeatedIntex(b, []byte("\r\n"), 2)]
		l2 := b[byteOps.RepeatedIntex(b, []byte("\r\n"), 2)+2 : byteOps.RepeatedIntex(b, []byte("\r\n"), 3)]

		if len(l0) >= 0 && byteOps.EndingOf(genBoundary(bou)[2:], l0) && (sufficientType(l1) == insufficient || sufficientType(l1) == sufficient) {

			resL = append(l0, []byte("\r\n")...)

			if sufficientType(l1) == insufficient {

				resL = append(resL, l1...)
				resL = append(resL, []byte("\r\n")...)

				if regexpops.IsCTFull(l2) {

					resL = append(resL, l2...)
					resL = append(resL, []byte("\r\n\r\n")...)

					return resL, fmt.Errorf("\"%s\" is %w", resL, errHeaderEnding)
				}
			}

		}
		if sufficientType(l0) == insufficient &&
			regexpops.IsCTFull(l1) {

			resL = append(resL, l0...)
			resL = append(resL, []byte("\r\n")...)
			resL = append(resL, l1...)
			resL = append(resL, []byte("\r\n\r\n")...)

			return resL, nil
		}
		if len(l0) == 0 { // on ending part is impossible, on beginning part: CRLF + CDsuf + 2*CRLF + rand || CRLF + CDinsuf + CRLF + CT + 2*CRLF + rand || CRLF + CT + 2*CRLF + rand || CRLF + rand // 2*CRLF + rand

			resL = append(l0, []byte("\r\n")...)

			if len(l1) == 0 {

				resL = append(resL, []byte("\r\n")...)

				return resL, fmt.Errorf("\"%s\" is %w", resL, errHeaderEnding)
			}

			if sufficientType(l1) == sufficient {

				resL = append(resL, l1...)
				resL = append(resL, []byte("\r\n\r\n")...)

				return resL, fmt.Errorf("\"%s\" is %w", resL, errHeaderEnding)
			}
			if sufficientType(l1) == insufficient {

				resL = append(resL, l1...)
				resL = append(resL, []byte("\r\n")...)

				if regexpops.IsCTFull(l2) {

					resL = append(resL, l2...)
					resL = append(resL, []byte("\r\n\r\n")...)

					return resL, fmt.Errorf("\"%s\" is %w", resL, errHeaderEnding)
				}
			}
			if regexpops.IsCTFull(l1) {

				resL = append(resL, l1...)
				resL = append(resL, []byte("\r\n\r\n")...)

				return resL, fmt.Errorf("\"%s\" is %w", resL, errHeaderEnding)
			}
			return resL, fmt.Errorf("\"%s\" is %w", resL, errHeaderEnding)
		}
		if len(l1) == 0 { // on ending part CDsuf + 2*CRLF + rand, on beginning part <-CDsuf + 2*CRLF + rand || <-CT + 2 * CRLF + rand

			resL = append(l0, []byte("\r\n\r\n")...)

			return resL, fmt.Errorf("\"%s\" is %w", resL, errHeaderEnding)
		}
		if len(l2) == 0 { // on ending part CDinsuf + CRLF + CT + 2*CRLF + rand, on beginning part CRLF + CDsuf + 2*CRLF = rand || <-Bound + CRLF + CDsuf + 2*CRLF = rand || <-CDinsuf + CRLF + CT + 2*CRLF

			resL = append(l0, []byte("\r\n")...)
			resL = append(resL, l1...)
			resL = append(resL, []byte("\r\n\r\n")...)

			return resL, fmt.Errorf("\"%s\" is %w", resL, errHeaderEnding)
		}
		return nil, errHeaderNotFound

	}
}

// sufficiency determines whether b is header for string data or for file data
func sufficientType(b []byte) sufficiency {

	r0 := regexp.MustCompile(`^Content-Disposition: form-data; name="[a-zA-zа-яА-Я0-9_.-:@#%^&\$\+\!\*\(\[\{\)\]\}]+"$`)
	r1 := regexp.MustCompile(`^Content-Disposition: form-data; name="[a-zA-zа-яА-Я0-9_.-:@#%^&\$\+\!\*\(\[\{\)\]\}]+"; filename="[a-zA-zа-яА-Я0-9_.-:@#%^&\$\+\!\*\(\[\{\)\]\}]+"$`)

	if r0.Match(b) {

		return sufficient
	}
	if r1.Match(b) {

		return insufficient
	}

	return incomplete
}

// isLastBoundaryPart returns true if b is ending part of last Boundary.
// Tested in dataHandler_test.go
func isLastBoundaryPart(b []byte, bou Boundary) bool {

	lenb, suffix := len(b), make([]byte, 0)

	i, lastSymbol := lenb, b[lenb-1]

	for i >= 1 {
		if i == 1 {

			return true
		}

		if i > 1 && b[i-1] != lastSymbol {

			break
		}

		i--
	}

	suffix = b[i:]
	rootLen := lenb - len(suffix)

	if rootLen < lenb && bytes.Contains(genBoundary(bou), b[:rootLen]) {

		return true
	}

	return false
}

// genBoundary generates byte slice based on given Boundary struct.
// Tested in dataHandler_test.go
func genBoundary(bou Boundary) []byte {

	Boundary := make([]byte, 0)

	Boundary = append(Boundary, []byte("\r\n")...)
	Boundary = append(Boundary, bou.Prefix...)
	Boundary = append(Boundary, bou.Root...)

	return Boundary
}

// getFoFi returns formname and filename found in b
func getFoFi(b []byte) (string, string) {

	fo, fi, foPre, fiPre := "", "", []byte(" name=\""), []byte(" filename=\"")

	if len(b) > 0 {

		if bytes.Contains(b, foPre) {

			fo = string(b[bytes.Index(b, foPre)+len(foPre) : byteOps.FindNext(b, []byte("\""), bytes.Index(b, foPre)+len(foPre))])

			if bytes.Contains(b, fiPre) {

				fi = string(b[bytes.Index(b, fiPre)+len(fiPre) : byteOps.FindNext(b, []byte("\""), bytes.Index(b, fiPre)+len(fiPre))])
			}
		}
	}

	return fo, fi
}

func getHeaderEndingIndex(body, header []byte) int {

	if len(body) == 0 || len(header) == 0 {

		return -1
	}
	i := 1

	for i < len(body) {

		if !bytes.Contains(header, body[:i]) {

			if i < 3 {

				return -1

			} else {

				return i - 1
			}

		}

		i++
	}

	return -1
}
