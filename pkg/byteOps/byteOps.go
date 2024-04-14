package byteOps

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
