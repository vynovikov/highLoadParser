package repository

import "errors"

var (
	errHeaderEnding   error = errors.New("ending of header")
	errHeaderNotFull  error = errors.New("header is not full")
	errHeaderNotFound error = errors.New("header is not found")
)

type TransferUnit struct {
}

type Presence struct {
}
