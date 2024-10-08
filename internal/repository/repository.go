package repository

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"

	"github.com/vynovikov/highLoadParser/internal/dataHandler"
	"github.com/vynovikov/highLoadParser/internal/logger"
	"github.com/vynovikov/highLoadParser/pkg/byteOps"
	regexpops "github.com/vynovikov/highLoadParser/pkg/regexpOps"
)

type ParserRepository interface {
	Register(RepositoryDTO) (dataHandler.ProducerUnit, error)
}

type repositoryStruct struct {
	dataHandler dataHandler.DataHandler
}

func NewParserRepository(dh dataHandler.DataHandler) *repositoryStruct {
	return &repositoryStruct{
		dataHandler: dh,
	}
}

func (r *repositoryStruct) Register(d RepositoryDTO) (dataHandler.ProducerUnit, error) {

	var (
		//err   error
		resTT dataHandler.ProducerUnit
	)

	switch {

	case d.B() == 0:

		key := newKeyDetailed(d)

		val, err := newValue(d)
		if err != nil &&
			!errors.Is(err, errHeaderNotFull) &&
			!errors.Is(err, errHeaderEnding) {

			return resTT, err
		}

		err = r.dataHandler.Set(key, val)
		if err != nil {

			logger.L.Infof("in repository.Register unable to set %s %d: %v\n", d.TS(), d.Part(), err)
		}

		/*	resTT, err = r.dataHandler.Create(d, bou)
				if err != nil {

					logger.L.Infof("in repository.Register unable to create %s %d: %v\n", d.TS(), d.Part(), err)
				}


			case d.B() == 1:

				bou := dataHandler.Boundary{}

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
		*/
	}

	return resTT, nil
}

func newValue(d RepositoryDTO) (dataHandler.Value, error) {

	headerB, body, _ := make([]byte, 0, maxHeaderLimit), d.Body(), 0

	lengh := len(body)

	if lengh > maxHeaderLimit {

		headerB = append(headerB, d.Body()[:maxHeaderLimit]...)

	} else {

		headerB = append(headerB, d.Body()...)
	}

	exactHeaderBytes, err := getHeaderLines(headerB, d.Bou())
	if err != nil {

		if errors.Is(err, errHeaderNotFull) ||
			errors.Is(err, errHeaderEnding) {

			return dataHandler.Value{
				E: d.E(),
				H: dataHandler.HeaderData{
					HeaderBytes: exactHeaderBytes,
				},
			}, err
		}

		return dataHandler.Value{}, err
	}

	fo, fi := getFoFi(exactHeaderBytes)

	return dataHandler.Value{
		E: d.E(),
		H: dataHandler.HeaderData{
			FormName:    fo,
			FileName:    fi,
			HeaderBytes: exactHeaderBytes,
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

func newKeyDetailed(d RepositoryDTO) dataHandler.KeyDetailed {

	return dataHandler.KeyDetailed{
		Ts:   d.TS(),
		Part: d.Part(),
	}
}

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

func genBoundary(bou Boundary) []byte {

	Boundary := make([]byte, 0)

	Boundary = append(Boundary, []byte("\r\n")...)
	Boundary = append(Boundary, bou.Prefix...)
	Boundary = append(Boundary, bou.Root...)

	return Boundary
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
// Tested in repository_test.go
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
