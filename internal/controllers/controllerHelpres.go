package controllers

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/vynovikov/highLoadParser/pkg/byteOps"
)

// AnalyzeHeader returns first 512 bytes of connection and boundary if found
func analyzeHeader(conn net.Conn) (boundary, []byte, error) {
	header := make([]byte, 512)
	conn.SetReadDeadline(time.Now().Add(time.Millisecond * 15)) // tls handshake requires at least 9 ms timeout

	n, err := io.ReadFull(conn, header)
	if err != nil &&
		(err != io.EOF && err != io.ErrUnexpectedEOF) ||
		(!os.IsTimeout(err) && n == 0) {
		return boundary{}, make([]byte, 0), err
	}
	//logger.L.Infof("in repo.AnalyzeHeader got from request %q\n", header)
	if n < len(header) {
		header = header[:n]
	}
	bou := findBoundary(header)
	if bytes.Contains(header, []byte("100-continue")) {
		return bou, header, fmt.Errorf("in repo.AnalyzeHeader expected 100-continue")
	}
	return bou, header, err
}

// AnalyzeBits returns result of reading 1024 bytes from connection
func analyzeBits(conn net.Conn, i, p int, h []byte, errFirst error) (parserControllerBody, error) {
	pcb, ending := newParserControllerBody(i), make([]byte, 0)
	if p == 0 &&
		(errFirst == nil || errFirst != nil && !strings.Contains(errFirst.Error(), "100-continue")) {

		lenh := len(h)
		if lenh < 512 {
			ending = make([]byte, 1024-lenh)
		} else {
			ending = make([]byte, 512)
		}

		pcb.B = h

		if lenh < 512 {
			return pcb, io.EOF
		}
		conn.SetReadDeadline(time.Now().Add(time.Millisecond * 1))
		n, err := io.ReadFull(conn, ending)

		if err != nil {
			if err != io.EOF && err != io.ErrUnexpectedEOF && !os.IsTimeout(err) {
				return pcb, err
			}
			// EOF
			if n > 0 && n <= len(ending) {

				ending = ending[:n]
				pcb.B = append(pcb.B, ending...)

				return pcb, err
			}

		}
		if n > 0 && n < len(ending) {
			ending = ending[:n]
		}
		pcb.B = append(pcb.B, ending...)

		return pcb, nil
	}
	conn.SetReadDeadline(time.Now().Add(time.Millisecond * 1))
	n, err := io.ReadFull(conn, pcb.B)
	if err != nil {

		if err != io.EOF && err != io.ErrUnexpectedEOF && !os.IsTimeout(err) {
			return pcb, err
		}
		// EOF
		if n == 0 {

			return newParserControllerBody(0), fmt.Errorf("in repo.AnalyzeBits request part %d is empty", p)
		}
		if n > 0 && n <= len(pcb.B) {

			pcb.B = pcb.B[:n]

			if errFirst != nil && strings.Contains(errFirst.Error(), "100-continue") {

				pcb.B = pcb.B[len(h):]
			}

			return pcb, err
		}

	}

	return pcb, nil

}

// Respond responds to connection with successful code
func respondOK(conn net.Conn) {
	body := "200 OK"
	doRespond(conn, body)

}

func respondContinue(conn net.Conn) {
	body := "100 Continue"
	doRespond(conn, body)
}

func doRespond(conn net.Conn, body string) {
	//logger.L.Infof("in repo.doRespond responding %q\n", body)
	fmt.Fprintf(conn, "HTTP/1.1 %s\r\nContent-Length: %d\r\nContent-Type: text/html\r\n\r\n%s", body, len(body), body)
}

func ReadFirst(conn net.Conn, n int) ([]byte, error) {
	firstN := make([]byte, n, n)

	conn.SetReadDeadline(time.Now().Add(time.Millisecond * 15))
	_, err := io.ReadFull(conn, firstN)
	if err != nil {
		return firstN, err
	}

	return firstN, nil
}

// FinsBoundary returns Boundary found in b
// Tested in byteOps_test.go
func findBoundary(b []byte) boundary {

	bPrefix, bRoot, bSuffix := make([]byte, 0, 2), make([]byte, 0, 48), make([]byte, 0, 2)

	if bytes.Contains(b, []byte(boundaryField)) {

		startIndex := bytes.Index(b, []byte(boundaryField)) + len(boundaryField)

		bRoot = byteOps.LineRightLimit(b, startIndex, 70)

		bPrefix = []byte("--")
	}
	return boundary{
		prefix: bPrefix,
		root:   bRoot,
		suffix: bSuffix,
	}

}

// GetLineWithCRLFLeft returns CRLF and succeeding line before given index.
// If line ends with CR (or CRLF) and contains boundary, returns CRLF + line + CR (CRLF).
// Tested in byteOps_test.go
func GetLineWithCRLFLeft(b []byte, fromIndex, limit int, bou boundary) []byte {

	l, lenb, c, n := make([]byte, 0), len(b), 0, 0

	if lenb < 1 {
		return l
	}
	if fromIndex > lenb-1 {
		fromIndex = lenb - 1
	}

	if fromIndex < limit {
		c = 0
	} else {
		c = fromIndex - limit
	}
	for i := fromIndex; i > c; i-- {
		if n == 0 &&
			(i == lenb-1 && b[i] == 13 ||
				i == lenb-2 && b[i] == 13 && b[i+1] == 10) &&
			(i >= 14 && ContainsBouEnding(b[i-14:i], bou)) {
			n++
			continue
		} else if i == lenb-1 && b[i] == 13 ||
			i == lenb-2 && b[i] == 13 && b[i+1] == 10 {

			return b[i:]
		}
		if b[i] == 13 && b[i+1] == 10 {
			return b[i:]
		}
	}
	return b
}

// ContainsBouEnding returns true if b contains boundary ending.
// Tested in byteOps_test.go
func ContainsBouEnding(b []byte, bou boundary) bool {
	n, boundary := 0, getBoundary(bou)
	for i := 0; i < len(b); i++ {
		if !bytes.Contains(boundary, b[:i]) && n > 4 {
			return true
		}
		if !bytes.Contains(boundary, b[:i]) {
			return false

		}
		n++
	}
	return true
}

// IsLastBoundary returns true if p + n form last boundary
func IsLastBoundary(p, n []byte, bou boundary) bool {
	realBoundary := getBoundary(bou)
	combined := append(p, n...)
	if len(combined) > len(realBoundary) &&
		(len(combined) > len(realBoundary)+1 && !bytes.Contains(combined[len(realBoundary):len(realBoundary)+2], []byte("\r\n")) ||
			len(combined) == len(realBoundary) && bytes.Contains(combined, realBoundary)) {
		return true
	}

	return false
}
