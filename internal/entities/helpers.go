package entities

import "bytes"

func GetBoundary(bou Boundary) []byte {

	boundary := make([]byte, 0)
	boundary = append(boundary, []byte("\r\n")...)
	boundary = append(boundary, bou.Prefix...)
	boundary = append(boundary, bou.Root...)

	return boundary
}

// GetLineWithCRLFLeft returns CRLF and succeeding line before given index.
// If line ends with CR (or CRLF) and contains boundary, returns CRLF + line + CR (CRLF).
// Tested in byteOps_test.go
func GetLineWithCRLFLeft(b []byte, fromIndex, limit int, bou Boundary) []byte {

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
func ContainsBouEnding(b []byte, bou Boundary) bool {
	n, boundary := 0, GetBoundary(bou)
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
func IsLastBoundary(p, n []byte, bou Boundary) bool {
	realBoundary := GetBoundary(bou)
	combined := append(p, n...)
	if len(combined) > len(realBoundary) &&
		(len(combined) > len(realBoundary)+1 && !bytes.Contains(combined[len(realBoundary):len(realBoundary)+2], []byte("\r\n")) ||
			len(combined) == len(realBoundary) && bytes.Contains(combined, realBoundary)) {
		return true
	}

	return false
}
