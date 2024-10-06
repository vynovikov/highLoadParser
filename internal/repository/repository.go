package repository

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/vynovikov/highLoadParser/internal/dataHandler"
	"github.com/vynovikov/highLoadParser/internal/logger"
	"github.com/vynovikov/highLoadParser/pkg/byteOps"
	regexpops "github.com/vynovikov/highLoadParser/pkg/regexpOps"
)

type ParserRepository interface {
	Register(dataHandler.DataHandlerDTO, dataHandler.Boundary) (dataHandler.ProducerUnit, error)
}

type repositoryStruct struct {
	dataHandler dataHandler.DataHandler
}

func NewParserRepository(dh dataHandler.DataHandler) *repositoryStruct {
	return &repositoryStruct{
		dataHandler: dh,
	}
}

func (r *repositoryStruct) Register(dto dataHandler.DataHandlerDTO, bou dataHandler.Boundary) (dataHandler.ProducerUnit, error) {

	var (
		err   error
		resTT dataHandler.ProducerUnit
	)

	d := dataHandler.NewDataHandlerUnit(dto)

	switch {

	case d.B() == 0:

		val, err := newValue(d, bou)
		if err != nil &&
			!errors.Is(err, errHeaderNotFull) &&
			!errors.Is(err, errHeaderEnding) {

			return resTT, err
		}

		err := r.dataHandler.Set(key, value)
		if err != nil {

			logger.L.Infof("in repository.Register unable to set %s %d: %v\n", d.TS(), d.Part(), err)
		}

	/*	resTT, err = r.dataHandler.Create(d, bou)
		if err != nil {

			logger.L.Infof("in repository.Register unable to create %s %d: %v\n", d.TS(), d.Part(), err)
		}
	*/
	case d.B() == 1:

		resTT, err = r.dataHandler.Updade(d, bou)
		if err != nil {

			logger.L.Infof("in repository.Register unable to update %s %d: %v\n", d.TS(), d.Part(), err)
		}
	}

	if d.Last() {

		err := r.dataHandler.Delete(d.TS())

		if err != nil {

			logger.L.Infof("in repository.Register unable to delete %s %v\n", d.TS(), err)
		}
	}

	return resTT, nil
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
