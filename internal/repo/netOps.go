package repo

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

// AnalyzeHeader returns first 512 bytes of connection and boundary if found
func AnalyzeHeader(conn net.Conn) (Boundary, []byte, error) {
	header := make([]byte, 512)
	conn.SetReadDeadline(time.Now().Add(time.Millisecond * 15)) // tls handshake requires at least 9 ms timeout
	n, err := io.ReadFull(conn, header)
	if err != nil &&
		(err != io.EOF && err != io.ErrUnexpectedEOF && !os.IsTimeout(err)) {
		return Boundary{}, header, err
	}
	//logger.L.Infof("in repo.AnalyzeHeader got from request %q\n", header)
	if n < len(header) {
		header = header[:n]
	}
	bou := FindBoundary(header)
	if bytes.Contains(header, []byte("100-continue")) {
		return bou, header, fmt.Errorf("in repo.AnalyzeHeader expected 100-continue")
	}
	return bou, header, err
}

// AnalyzeBits returns result of reading 1024 bytes from connection
func AnalyzeBits(conn net.Conn, i, p int, h []byte, errFirst error) (ReceiverBody, error) {
	rb, ending := NewReceiverBody(i), make([]byte, 0)
	if p == 0 &&
		(errFirst == nil || errFirst != nil && !strings.Contains(errFirst.Error(), "100-continue")) {

		lenh := len(h)
		if lenh < 512 {
			ending = make([]byte, 1024-lenh)
		} else {
			ending = make([]byte, 512)
		}

		rb.B = h

		if lenh < 512 {
			return rb, io.EOF
		}
		conn.SetReadDeadline(time.Now().Add(time.Millisecond * 1))
		n, err := io.ReadFull(conn, ending)

		if err != nil {
			if err != io.EOF && err != io.ErrUnexpectedEOF && !os.IsTimeout(err) {
				return rb, err
			}
			// EOF
			if n > 0 && n <= len(ending) {

				ending = ending[:n]
				rb.B = append(rb.B, ending...)

				return rb, err
			}

		}
		if n > 0 && n < len(ending) {
			ending = ending[:n]
		}
		rb.B = append(rb.B, ending...)

		return rb, nil
	}
	conn.SetReadDeadline(time.Now().Add(time.Millisecond * 1))
	n, err := io.ReadFull(conn, rb.B)
	if err != nil {

		if err != io.EOF && err != io.ErrUnexpectedEOF && !os.IsTimeout(err) {
			return rb, err
		}
		// EOF
		if n == 0 {

			return NewReceiverBody(0), fmt.Errorf("in repo.AnalyzeBits request part %d is empty", p)
		}
		if n > 0 && n <= len(rb.B) {

			rb.B = rb.B[:n]

			return rb, err
		}

	}

	return rb, nil

}

// Respond responds to connection with successful code
func Respond(conn net.Conn) {

	body := "200 OK"
	doRespond(conn, body)
}

func RespondContinue(conn net.Conn) {

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
