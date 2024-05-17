package byteOps

import "bytes"

// LineRightLimit returns byte slice to the right of fromIndex in b. Stops when CR found of limit exceeded
// Tested in byteOps_test.go
func LineRightLimit(b []byte, fromIndex, limit int) []byte {
	bb := make([]byte, 0)

	if fromIndex < 0 {
		return nil
	}

	for i := fromIndex; b[i] != 13 && i < fromIndex+limit; i++ {
		bb = append(bb, b[i])
	}
	if len(bb) == limit {
		return nil
	}

	return bb
}

// BeginningEqual returns true if first and second slices have equal characters in the same positions.
// Tested in byteOps_test.go
func BeginningEqual(s1, s2 []byte) bool {
	if len(s1) > len(s2) {
		s1 = s1[:len(s2)]
	} else {
		s2 = s2[:len(s1)]
	}
	for i, v := range s2 {
		if s1[i] != v {
			return false
		}
	}
	return true
}

// RepeatedIntex returns index of not first occurence of occ in byte slice.
// Tested in byteOps_test.go
func RepeatedIntex(b, occ []byte, i int) int {
	index, n := 0, 0

	for n < i {
		n++
		indexN := bytes.Index(b, occ)
		if n == 1 {
			index += indexN
		} else {
			index += indexN + len(occ)
		}
		cutted := indexN + len(occ)

		b = b[cutted:]

	}
	return index
}

// EndingOf returns true if first slice contains second and second slice is the ending of the first.
// Tested in byteOps_test.go
func EndingOf(long, short []byte) bool {
	longtLE, shortLE, lenLong, lenShort := byte(0), byte(0), len(long), len(short)
	if lenShort < 1 {
		return true
	}
	if lenLong < 1 {
		return false
	}
	shortLE = short[lenShort-1]
	longtLE = long[lenLong-1]
	if longtLE != shortLE {
		return false
	}
	for i := lenShort - 1; i > -1; i-- {
		if short[i] != long[lenLong-lenShort+i] {
			return false
		}
	}

	return true
}

// FindNext returns index of occ next to fromIndex
func FindNext(b, occ []byte, fromIndex int) int {

	return bytes.Index(b[fromIndex:], occ) + fromIndex
}
