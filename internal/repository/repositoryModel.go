package repository

import "errors"

type (
	Disposition int
	sufficiency int
)

var (
	errHeaderEnding   error = errors.New("ending of header")
	errHeaderNotFull  error = errors.New("header is not full")
	errHeaderNotFound error = errors.New("header is not found")
)

const (
	CONTENT_DISPOSITION             = "Content-Disposition"
	False               Disposition = iota
	True
	Probably
	incomplete sufficiency = iota
	sufficient
	insufficient
	sep            = "\r\n"
	maxLineLimit   = 100
	maxHeaderLimit = 210
)

type TransferUnit struct {
}

type Presence struct {
}

type RepositoryDTO interface {
	Part() int
	TS() string
	Body() []byte
	SetBody([]byte)
	B() int
	E() int
	Last() bool
	IsSub() bool
}

type Boundary struct {
	Prefix []byte
	Root   []byte
	Suffix []byte
}

type RepositoryUnit struct {
	R_part     int
	R_ts       string
	R_body     []byte
	R_b        int
	R_e        int
	R_last     bool
	R_isSub    bool
	R_boundary Boundary
}
