package controllers

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/vynovikov/highLoadParser/internal/entities"
)

// AnalyzeHeader returns first 512 bytes of connection and boundary if found
func analyzeHeader(conn net.Conn) (entities.Boundary, []byte, error) {

	header := make([]byte, 512)

	conn.SetReadDeadline(time.Now().Add(time.Millisecond * 15)) // tls handshake requires at least 9 ms timeout

	n, err := io.ReadFull(conn, header)
	if err != nil {
		if err, ok := err.(net.Error); ok && !err.Timeout() {
			return entities.Boundary{}, make([]byte, 0), err
		}
	}
	if n < len(header) {
		header = header[:n]
	}

	bou := entities.FindBoundary(header)
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
	defer conn.Close()

	// Write HTTP status line and headers
	fmt.Fprintf(conn, "HTTP/1.1 %s\r\n", body)
	fmt.Fprint(conn, "Content-Type: text/plain\r\n")
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(body)+len("\r\n")) // length of the response body
	fmt.Fprint(conn, "\r\n")                                           // end of headers

	fmt.Fprintf(conn, "%s\r\n", body)

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
