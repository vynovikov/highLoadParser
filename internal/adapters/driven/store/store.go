package store

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"workspaces/postParser/internal/repo"
)

type StoreStruct struct {
	R map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue
	B map[repo.AppStoreKeyGeneral][]repo.DataPiece
	C map[repo.AppStoreKeyGeneral]repo.Counter
	L sync.RWMutex
}

func NewStore() *StoreStruct {
	return &StoreStruct{
		R: make(map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue),
		B: make(map[repo.AppStoreKeyGeneral][]repo.DataPiece),
		C: make(map[repo.AppStoreKeyGeneral]repo.Counter),
		L: sync.RWMutex{},
	}
}

type Store interface {
	IsBufferEmpty() bool
	//DetailedKey(string) repo.AppStoreKeyDetailed
	BufferAdd(repo.DataPiece)
	Register(repo.DataPiece, repo.Boundary) (repo.AppDistributorUnit, error)
	RegisterBuffer(repo.AppStoreKeyGeneral, repo.Boundary) ([]repo.AppDistributorUnit, []error)
	Inc(repo.AppStoreKeyGeneral, int)
	Dec(repo.DataPiece) (repo.Order, error)
	Counter(repo.AppStoreKeyGeneral) int
	//Reset(repo.AppStoreKeyGeneral)
	Presence(repo.DataPiece) (repo.Presense, error)
	Update(repo.DetailedRecord)
	Delete(repo.AppStoreKeyGeneral)
	Act(repo.DataPiece, repo.StoreChange)
	Unblock(repo.AppStoreKeyGeneral)
}

func (s *StoreStruct) IsBufferEmpty() bool {
	return len(s.B) == 0
}

func (s *StoreStruct) HasGeneral(ts string) bool {
	/*
		askg := repo.NewAppStoreKeyGeneral(ts)
		if _, ok := s.R[askg]; ok {
			return true
		}
	*/
	return false
}

func (s *StoreStruct) DetailedPart(ts string) (int, error) {
	/*parts := make([]int, 0)

	askg := repo.NewAppStoreKeyGeneral(ts)
	d, ok := s.R[askg]
	//logger.L.Infof("in store.DetailedPart d = %c, ok %t", d, ok)
	if ok {
		for i, _ := range d {
			parts = append(parts, i.Part)
		}

		if len(parts) > 0 {
			return parts[0], nil
		}
	}
	*/
	return -1, errors.New("not found")

}

func (s *StoreStruct) DetailedKey(ts string) repo.AppStoreKeyDetailed {
	iMax := repo.AppStoreKeyDetailed{}
	/*
		askg := repo.NewAppStoreKeyGeneral(ts)
		d, ok := s.R[askg]

		if ok {
			for i := range d {
				if i.Part > iMax.Part {
					iMax = i
				}
			}

		}
	*/
	return iMax
}

func (s *StoreStruct) Register(d repo.DataPiece, bou repo.Boundary) (repo.AppDistributorUnit, error) {

	//logger.L.Infof("store.Register invoked with d header: %v, body: %q, while s.R: %v, s.C: %v\n", d.GetHeader(), d.GetBody(0), s.R, s.C)

	var (
		du        repo.AppDistributorUnit
		askdParts []int
	)
	if s.B == nil { // delete after testing
		s.B = make(map[repo.AppStoreKeyGeneral][]repo.DataPiece)
	}

	askg := repo.NewAppStoreKeyGeneralFromDataPiece(d)
	//logger.L.Infof("in store.Register TS: %q, askg: %q\n", d.TS(), askg)
	askd := repo.NewAppStoreKeyDetailed(d)
	//logger.L.Infof("in store.Register TS: %q, askd: %q\n", d.TS(), askd)

	switch {

	case d.B() == repo.True: //DataPiece needs beginning

		if m1, ok := s.R[askg]; ok { // General key met

			if m2, ok := m1[askd]; ok { // Detailed key met

				//logger.L.Infof("in store.Register specific map is %v\n", m2)
				//logger.L.Infof("in store.Register in dataPiece TS: %s, Part = %d, E() = %d, body: %q\n", d.TS(), d.Part(), d.E(), d.GetBody(0))

				m3f := m2[false]
				//logger.L.Infof("in store.Register right place\n")
				if m3f.D.FormName == "" { // Disposition is not filled
					header, err := d.H(bou)
					//logger.L.Warnf("in store.Register header: %q, error %v\n", header, err)

					if m3t, ok := m2[true]; ok { // true branch is present
						if strings.Contains(err.Error(), "is ending part") {

							m3f.D.H = append(m3f.D.H, m3t.D.H...)
							m3f.D.H = append(m3f.D.H, header...)
							m3f.D.FormName, m3f.D.FileName = repo.GetFoFi(m3f.D.H)
							m3f.E = d.E()

							d.BodyCut(len(header))

							//logger.L.Infof("in store.Register m3t.D.H: %q, m3t.D.FormName: %q, m3t.D.FileName: %q\n", m3t.D.H, m3t.D.FormName, m3t.D.FileName)

							vvv := map[bool]repo.AppStoreValue{}
							vvv[false] = m3f
							delete(s.R[askg], askd)
							s.R[askg][askd.IncPart()] = vvv

							return repo.NewDistributorUnitStream(m3f, d, repo.Message{PreAction: repo.Continue, PostAction: repo.Continue}), fmt.Errorf("in store.Register new dataPiece group with TS \"%s\" and Part = %d", askg.TS, m3f.B.Part)
						}

					}
					// no true branch
					if strings.Contains(err.Error(), "is ending part") {

						m3f.D.H = append(m3f.D.H, header...)
						m3f.D.FormName, m3f.D.FileName = repo.GetFoFi(m3f.D.H)
						m3f.E = d.E()

						d.BodyCut(len(header))

						vvv := map[bool]repo.AppStoreValue{}
						vvv[false] = m3f

						//logger.L.Infof("in store.Register s.M: %v\n", s.M)
						delete(s.R[askg], askd)

						s.R[askg][askd.IncPart()] = vvv

						//logger.L.Infof("in store.Register s.M: %v\n", s.M)
						//logger.L.Infof("in store.Register adu header: %v, body: %q\n", adu.H, adu.B.B)
						return repo.NewDistributorUnitStream(m3f, d, repo.Message{PreAction: repo.Continue, PostAction: repo.Continue}), fmt.Errorf("in store.Register new dataPiece group with TS \"%s\" and Part = %d", askg.TS, m3f.B.Part)
					}
				}
				// Disposition is filled
				if m3t, ok := m2[true]; ok { // true branch is present

					if m3f.E == repo.True { // false branch needs next part -- still waiting for dataPiece with d.E == repo.Probably which is current one

						//logger.L.Infof("in store.Register is this case with m3f: %q\n", m3f)

						m3f.E = d.E()

						vv := map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{askd.IncPart(): {false: m3f, true: m3t}}
						delete(s.R[askg], askd)
						s.R[askg] = vv
						//logger.L.Warnf("in store.Register s.R: %v\n", s.R)

						return repo.NewDistributorUnitStream(m3f, d, repo.Message{PreAction: repo.Continue, PostAction: repo.Continue}), nil

					}

					header, err := d.H(bou)
					//logger.L.Warnf("in store.Register in dataPiece TS: %s, Part = %d, E() = %d -> header: %q, error: %v \n", d.TS(), d.Part(), d.E(), header, err)

					if err != nil {
						if strings.Contains(err.Error(), "is ending part") { // header in dataPiece's body present

							if repo.IsLastBoundary(m3t.D.H, header, bou) {
								adu := repo.AppDistributorUnit{H: repo.AppDistributorHeader{T: repo.ClientStream, S: repo.StreamData{SK: repo.StreamKey{TS: d.TS(), Part: m3t.B.Part}, M: repo.Message{S: "finish"}}}}
								delete(s.R, askg)

								return adu, fmt.Errorf("in store.Register dataPiece group with TS \"%s\" is finished", d.TS())
							}

							if repo.IsBoundary(m3t.D.H, header, bou) && bytes.Contains(header, []byte("Content-Disposition")) {

								//logger.L.Infoln("in store.Register this case")

								m3f.D.H = header[bytes.Index(header, []byte("Content-Disposition")):]
								m3f.D.FormName, m3f.D.FileName = repo.GetFoFi(header)
								m3f.B = m3t.B
								m3f.E = d.E()

								//body:=d.GetBody(0)[len(header):]
								d.BodyCut(len(header))

								vv := map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{askd.IncPart(): {false: m3f}}

								delete(s.R[askg], askd)

								s.R[askg] = vv

								adu := repo.AppDistributorUnit{}

								switch d.E() {
								case repo.True:
									adu = repo.NewDistributorUnitStream(m3f, d, repo.Message{PreAction: repo.StopLast, PostAction: repo.Continue})
								case repo.False:
									//logger.L.Infof("in store.Register dataPiece body: %q\n", d.GetBody(0))
									adu = repo.NewAppDistributorUnitUnaryComposed(m3f, d, repo.Message{PreAction: repo.StopLast, PostAction: repo.Close})
								}

								//logger.L.Warnf("in store.Register adu header %v, body %q\n", adu.H, adu.B.B)

								return adu, fmt.Errorf("in store.Register new dataPiece group with TS \"%s\" and Part = %d", askg.TS, m3f.B.Part)
							}

							d.Prepend(m3t.D.H)

							m3f.E = d.E()

							vv := map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{askd.IncPart(): {false: m3f}}

							delete(s.R[askg], askd)

							s.R[askg] = vv

							return repo.NewDistributorUnitStream(m3f, d, repo.Message{PreAction: repo.Continue, PostAction: repo.Continue}), nil

						}

						if strings.Contains(err.Error(), "no header found") {

							m3f.E = d.E()

							vv := map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue{askd.IncPart(): {false: m3f}}
							delete(s.R[askg], askd)
							s.R[askg] = vv

							d.Prepend(m3t.D.H)

							return repo.NewDistributorUnitStream(s.R[askg][askd.IncPart()][false], d, repo.Message{PreAction: repo.Continue, PostAction: repo.Continue}), nil

						}
						if strings.Contains(err.Error(), "is the last") {

							delete(s.R, askg)

							adu := repo.NewDistributorUnitStreamEmpty(m3f, d, repo.Message{PreAction: repo.Continue, PostAction: repo.Finish})

							return adu, fmt.Errorf("in store.Register dataPiece group with TS \"%s\" is finished", d.TS())
						}
					}
					if repo.IsLastBoundary(m3t.D.H, header, bou) {
						return repo.AppDistributorUnit{H: repo.AppDistributorHeader{T: repo.Tech, C: repo.CloseData{TS: d.TS()}}}, nil
					}
				}
				// no true branch
				vvv := map[bool]repo.AppStoreValue{}
				switch d.E() {

				case repo.Probably: // deteiled group can be closed or finished

					//logger.L.Warnf("in store.Register s.R before: %v\n", s.R)
					//logger.L.Warnf("in store.Register m1: %v, len(m1) = %d\n", m1, len(m1))

					if _, ok := m2[true]; ok { // level 3 map has true branch

						//logger.L.Infof("in store.Register got right place with previously saved header line %q and newly body %q\n", m3t.D.H[0], d.GetBody())
						//increment askd

					}
					m2f, ok1 := m1[askd.F()]
					m2t, ok2 := m1[askd.T()]
					if ok1 && ok2 { // there is another ASKD with same part which has true branch only
						//logger.L.Infof("in store.Register m1: %v, askd.T(): %v\n", m1, askd.T())
						m3t := m2t[true]

						//logger.L.Infof("in store.Register m3t: %v\n", m3t)
						m3f := m2f[false]
						m3f.E = d.E()
						//logger.L.Infof("in store.Register m3f: %v\n", m3f)

						vvv := map[bool]repo.AppStoreValue{}
						vvv[true] = m3t
						vvv[false] = m3f

						delete(s.R[askg], askd.T())
						delete(s.R[askg], askd.F())
						//logger.L.Warnf("in store.Register s.R after deleting: %v\n", s.R)
						s.R[askg][askd.IncPart()] = vvv
						//logger.L.Warnf("in store.Register s.R after adding: %v\n", s.R)
						return repo.NewDistributorUnitStream(m3f, d, repo.Message{PreAction: repo.Continue, PostAction: repo.Continue}), nil
						/*						for i,v:=range m1{
													if m3t,ok:=v[true];ok&&len(v)==1[

													]
												}
						*/
					}
					// no true branch

					m3f.E = d.E()

					delete(s.R[askg], askd)
					s.R[askg][askd] = map[bool]repo.AppStoreValue{false: m3f}

					return repo.NewDistributorUnitStream(m3f, d, repo.Message{PreAction: repo.Continue, PostAction: repo.Continue}), fmt.Errorf("in store.Register got double-meaning dataPiece")

				case repo.True: // detailed group shall be continued

					vvv[false] = m3f
					/*
						if m2, ok := s.R[askg][askd.IncPart()]; ok && len(m2) == 1 {
							if m3t, ok := m2[true]; ok {
								vvv[true] = m3t
							}
						}
					*/
					delete(s.R[askg], askd)
					//logger.L.Warnf("in store.Register after deleting s.R = %v\n", s.R)

					s.R[askg][askd.IncPart()] = vvv
					//logger.L.Warnf("in store.Register s.R = %v\n", s.R)

					adu := repo.NewDistributorUnitStream(m3f, d, repo.Message{PreAction: repo.Continue, PostAction: repo.Continue})
					//logger.L.Warnf("in store.Register adu header: %v, body: %q\n", adu.H, adu.B.B)

					return adu, nil

				case repo.False:

					//logger.L.Infof("in store.Register d.B(), d.E()==repo.False \n")

					if m3t, ok := m2[true]; ok {

						boundaryTrimmed := repo.GenBoundary(bou)[2:]
						lb := append(m3t.D.H, d.GetBody(repo.Min(d.LL(), repo.MaxHeaderLimit))...)

						if repo.LastBoundary(lb, boundaryTrimmed) {
							//logger.L.Infof("in store.Register got last boundary %q\n", lb)
							//s.R = make(map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue)
							//logger.L.Infof("in store.Register deleting askg %q\n", askg)
							delete(s.R, askg)
							//logger.L.Infof("in store.Register after deleting askg s.R: %v\n", s.R)
							return repo.NewDistributorUnitStreamEmpty(m3t, d, repo.Message{PreAction: repo.Continue, PostAction: repo.Close}), fmt.Errorf("in store.Register dataPiece group with TS \"%s\" is finished", d.TS())
						}

						delete(s.R[askg], askd)

						return repo.NewDistributorUnitStream(m3f, d, repo.Message{PreAction: repo.Continue, PostAction: repo.Close}), fmt.Errorf("in store.Register dataPiece group with TS %q and Part = %d is finished", d.TS(), m3f.B.Part)
					}
					//logger.L.Infof("in store.Register before deleting s.R %v, len(s.R[askg]) = %d\n", s.R[askg], len(s.R[askg]))

					delete(s.R[askg], askd)

					//logger.L.Infof("in store.Register after registering apu header %v, s.R: %v\n", d.GetHeader(), s.R)

					return repo.NewDistributorUnitStream(m3f, d, repo.Message{PreAction: repo.Continue, PostAction: repo.Close}), fmt.Errorf("in store.Register dataPiece group with TS %q and Part \"%d\" is finished", d.TS(), m3f.B.Part)
				}
			}
			// askd not met

			s.BufferAdd(d)

			for i := range s.R[askg] {
				//logger.L.Warnf("in store.Register ranging over s.R => askd = %v\n", i)

				askdParts = append(askdParts, i.SK.Part)

			}
			sort.Ints(askdParts)
			//logger.L.Warnf("in store.Register askdParts: %d\n", askdParts)
			if len(askdParts) == 1 {
				return du, fmt.Errorf("in store.Register dataPiece's Part for given TS \"%s\" should be \"%d\" but got \"%d\"", d.TS(), askdParts[0], d.Part())
			}
			if len(askdParts) > 1 {
				return du, fmt.Errorf("in store.Register dataPiece's Part for given TS \"%s\" should be one of \"%d\" but got \"%d\"", d.TS(), askdParts, d.Part())
			}

		}
		// askg not met

		s.BufferAdd(d)
		return du, fmt.Errorf("in store.Register dataPiece's TS %q is unknown", askg.TS)
	}
	return du, nil
}

func (s *StoreStruct) BufferAdd(d repo.DataPiece) {

	askg := repo.NewAppStoreKeyGeneralFromDataPiece(d)
	if s.B == nil {
		s.B = map[repo.AppStoreKeyGeneral][]repo.DataPiece{}
	}

	switch lenb := len(s.B[askg]); {
	case lenb == 0:
		s.B[askg] = make([]repo.DataPiece, 0)
		s.B[askg] = append(s.B[askg], d)

	case lenb == 1:

		e := s.B[askg][0]

		if Equal(d, e) { // avoiding dublicates
			return
		}

		if d.Part() >= e.Part() {
			//logger.L.Infof("in store.BufferAdd in case 1 trying to append dataPiece with body %q to buffer %v\n", d.GetBody(0), s.B[askg])
			s.B[askg] = append(s.B[askg], d)
		}
		if d.Part() < e.Part() {
			s.Prepend(s.B[askg], d)
		}

	default:

		for _, v := range s.B[askg] { // avoiding dublicates
			if Equal(d, v) {
				return
			}
		}

		f := s.B[askg][0]
		l := s.B[askg][len(s.B[askg])-1]

		if f.Part()+1 != l.Part() { //if first and last elements are not neighbours

			if d.Part() > l.Part() { //new element's Part is more than last element's
				s.B[askg] = append(s.B[askg], d)
			}
			if d.Part() < f.Part() { //new element's Part is less than first element's
				s.Prepend(s.B[askg], d)
			}
			if d.Part() == l.Part()-1 {

				lastIndex := len(s.B[askg]) - 1
				s.B[askg] = append(s.B[askg], d)

				Swap(s.B[askg], lastIndex, lastIndex+1)

			}
			if d.Part() == f.Part()+1 {

				s.Prepend(s.B[askg], d)

				Swap(s.B[askg], 0, 1)

			}
			if d.Part() < l.Part()-1 && //new element's Part is less than last element's
				d.Part() > f.Part()+1 { //new element's Part is more than first element's

				s.B[askg] = append(s.B[askg], d)
				b := s.B[askg]
				sort.SliceStable(s.B[askg], func(i int, j int) bool { return b[i].Part() < b[j].Part() })
			}

		}
		if f.Part()+1 == l.Part() { //if first and last element are neighbours
			if d.Part() > l.Part() { //new element's Part is more than last element's
				s.B[askg] = append(s.B[askg], d)
			}
			if d.Part() < f.Part() { //new element's Part is less than first element's
				s.Prepend(s.B[askg], d)
			}
		}
	}
	//logger.L.Infof("in store.BufferAdd buffer became %v\n", s.B[askg])
}

// Inserting new element to the beginning of slice, shifting the rest
func (s *StoreStruct) Prepend(tsb []repo.DataPiece, d repo.DataPiece) {

	askg := repo.NewAppStoreKeyGeneralFromDataPiece(d)

	if b, ok := s.B[askg]; ok {
		s.B[askg] = append([]repo.DataPiece{}, d)
		s.B[askg] = append(s.B[askg], b...)
		//logger.L.Infof("in store.Prepend len became %d\n", len(s.B[d.TS()]))
		//logger.L.Infof("in store.Prepend s.B[d.TS()][1] = %v\n", s.B[d.TS()][1])
		return
	}
	s.B[askg] = append([]repo.DataPiece{}, d) // something wrong case
}

// Swapping two slice elements
func Swap(s []repo.DataPiece, i, j int) {
	e := s[i]
	s[i] = s[j]
	s[j] = e
}

func (s *StoreStruct) CleanBuffer(ids repo.AppStoreBufferIDs) {
	s.L.Lock()
	defer s.L.Unlock()
	//logger.L.Infof("store.CleanBuffer invoked with s.B: %v, ids: %v, len = %d\n", s.B, ids, len(ids.I))

	marked := 0
	if len(ids.I) == 0 {
		return
	}

	askg := ids.ASKG
	// setting registerd dataPieces as empty
	for i := range ids.I {
		//logger.L.Infof("store.CleanBuffer s.B[askg]: %v, ids.I[i] = %d\n", s.B[askg], ids.I[i])
		s.B[askg][ids.I[i]] = &repo.AppPieceUnit{}
		marked++
	}
	if marked == len(s.B[askg]) {
		delete(s.B, askg)
		return
	}
	// sorting B putting empty ones to buffer beginning

	sort.SliceStable(s.B[askg], func(a, b int) bool {
		return s.B[askg][a].TS() < s.B[askg][b].TS()
	})
	//cutting off the buffer beginning
	for j := range s.B[askg] {
		//logger.L.Infof("store.CleanBuffer j = %d, s.B[ids.TS][j].TS(): %v\n", j, s.B[ids.TS][j].TS())
		if s.B[askg][j].TS() != "" {
			s.B[askg] = s.B[askg][j:]
			break
		}
		if j == len(s.B[askg])-1 {
			s.B[askg] = []repo.DataPiece{}
		}
	}
}

// Register buffer elements, considering they are sorted by Part
func (s *StoreStruct) RegisterBuffer(askg repo.AppStoreKeyGeneral, bou repo.Boundary) ([]repo.AppDistributorUnit, []error) {
	//logger.L.Infof("store.RegisterBuffer invoked, askg %v, s.R[askg]: %v, s.B[askg]: %v and s.C[askg]: %v\n", askg, s.R[askg], s.B[askg], s.C[askg])
	/*
		for i, v := range s.B[askg] {
			logger.L.Infof("in store.RegisterBuffer buffer element %d is %v, s.C[askg].Cur = %v\n", i, v.GetHeader(), s.C[askg].Cur)
		}
	*/
	ids, adus, errs, repeat, _ := repo.NewAppStoreBufferIDs(), make([]repo.AppDistributorUnit, 0), make([]error, 0), true, len(s.B[askg]) > 0 && !s.C[askg].Blocked
	//logger.L.Infof("in store.RegisterBuffer debug: %t\n", debug)
	/*p,_:=d.Part(),d.E()
	  logger.L.Infof("in store.RegisterBuffer after register ")
	*/

	if len(s.B) == 0 {
		//logger.L.Errorln("in store.RegisterBuffer buffer is empty, leaving")
		errs = append(errs, fmt.Errorf("in store.RegisterBuffer buffer has no elements"))
		return adus, errs
	}
	//askg := repo.NewAppStoreKeyGeneralFromDataPiece(d)
	if len(s.B[askg]) == 1 {
		if s.B[askg][0].B() == repo.False && s.B[askg][0].E() == repo.False {
			if s.C[askg].Cur == 1 && !s.C[askg].Blocked {
				switch s.C[askg].Max {
				case 1:
					adus = append(adus, repo.NewAppDistributorUnitUnary(s.B[askg][0], bou, repo.Message{PreAction: repo.Start, PostAction: repo.Finish}))
				default:
					adus = append(adus, repo.NewAppDistributorUnitUnary(s.B[askg][0], bou, repo.Message{PreAction: repo.None, PostAction: repo.Finish}))
				}
				return adus, nil

			}
		}
	}
	if len(s.B[askg]) == 1 && s.C[askg].Cur == 1 {
		if s.C[askg].Blocked {
			errs = append(errs, fmt.Errorf("in store.RegisterBuffer buffer has single element and current counter == 1 and blocked"))
			return adus, errs
		}

	}

	for i, v := range s.B {
		//logger.L.Infof("in store.RegisterBuffer checking group of TS \"%s\"\n", i.TS)
		if i.TS != askg.TS {
			//logger.L.Errorf("in store.RegisterBuffer askg TS is \"%s\" != d.TS() \"%s\"\n", i.TS, d.TS())
			continue
		}

		for repeat {

			repeat = false

			for j, w := range v {

				askg, askd := repo.NewAppStoreKeyGeneralFromDataPiece(w), repo.NewAppStoreKeyDetailed(w)
				//logger.L.Infof("in store.RegisterBuffer j = %d w header %v body %q, counter %d, ids %d, s.B[askg] %v\n", j, w.GetHeader(), w.GetBody(0), s.Counter(askg), ids.GetIDs(askg), s.B[askg])

				if isIn(j, ids.GetIDs(askg)) { // index'v been registered already, skipping

					if len(ids.GetIDs(askg)) == len(s.B[askg]) {

						break
					}

					continue
				}

				if m1, ok := s.R[askg]; ok {
					//logger.L.Infoln("in store.RegisterBuffer askg matched")
					if _, ok := m1[askd]; ok {
						//logger.L.Infoln("in store.RegisterBuffer askd matched")
						//logger.L.Infof("in store.RegisterBuffer checking dataPiece with header %v, body %q, counter = %v\n", w.GetHeader(), w.GetBody(0), s.C)
						//logger.L.Infof("in store.RegisterBuffer dataPiece with header %v matched to R\n", w.GetHeader())
						if !s.C[askg].Blocked {
							adu, err := s.Register(w, bou)
							//logger.L.Infof("in store.RegisterBuffer after register adu header: %v, s.R: %v, counter: %v\n", adu.H, s.R, s.C)
							if err != nil {
								errs = append(errs, err)
							}
							if s.C[askg].Cur == 1 {
								if adu.H.T == repo.Unary {
									adu.H.U.M.PostAction = repo.Finish
								}
								if adu.H.T == repo.ClientStream {
									adu.H.S.M.PostAction = repo.Finish
								}
								//logger.L.Infof("in store.RegisterBuffer after setting finish adu header: %v\n", adu.H)
							}
							s.Dec(w)

							adus = append(adus, adu)
							ids.Add(repo.NewAppStoreBufferID(i, j))
							//logger.L.Infof("in store.RegisterBuffer ids became askg: %v, ids.I %d, len(ids.I) = %d\n", ids.ASKG, ids.I, len(ids.I))

							continue
						}
						if s.C[askg].Blocked && s.C[askg].Cur > 1 {
							adu, err := s.Register(w, bou)
							//logger.L.Infof("in store.RegisterBuffer after register adu header: %v, s.R: %v, counter: %v\n", adu.H, s.R, s.C)
							if err != nil {
								errs = append(errs, err)
							}

							s.Dec(w)

							//logger.L.Infof("in store.RegisterBuffer after counter check adu header: %v, body: %q\n", adu.H, adu.B.B)
							adus = append(adus, adu)
							//logger.L.Infof("in store.RegisterBuffer counter decremented and became %d\n", s.Counter(askg))
							repeat = true
							ids.Add(repo.NewAppStoreBufferID(i, j))
							//logger.L.Infof("in store.RegisterBuffer ids became askg: %v, ids %d\n", ids.ASKG, ids.I)
							continue

						}
						if s.C[askg].Blocked && s.C[askg].Cur == 1 {
							errs = append(errs, fmt.Errorf("in store.RegisterBuffer buffer has single element and current counter == 1 and fuse"))
							repeat = false
							break
						}

					}
				}
				repeat = false
			}
		}
		//logger.L.Infof("in store.RegisterBuffer checking group of TS \"%s\" is finiahed\n", i.TS)
	}
	//logger.L.Infof("in store.RegisterBuffer before clearing buffer ids: %v, len(ids.I) = %d, len(s.B[askg]) = %d\n", ids, len(ids.I), len(s.B[askg]))
	if len(ids.I) > 0 && len(s.B[askg]) > 0 {
		//logger.L.Infof("in store.RegisterBuffer trying to clear buffer ids: %v\n", ids)
		s.CleanBuffer(*ids)
	}

	//logger.L.Infof("in store.RegisterBuffer after all s.R: %v, s.B: %v, s.C: %v\n", s.R, s.B, s.C)
	return adus, errs
}

func (s *StoreStruct) Inc(askg repo.AppStoreKeyGeneral, n int) {
	s.L.Lock()
	defer s.L.Unlock()
	c := repo.NewCounter()

	if cm, ok := s.C[askg]; ok {
		cm.Max += n
		cm.Cur += n
		s.C[askg] = cm
		return
	}
	c.Max = n
	c.Cur = n
	c.Blocked = true
	s.C[askg] = c

}

// Decrements counter if possible and returns result and error
func (s *StoreStruct) Dec(d repo.DataPiece) (repo.Order, error) {
	//logger.L.Infof("store.Dec invoked with dataPiece header %v,body %q, s.C: %v\n", d.GetHeader(), d.GetBody(0), s.C)
	askg := repo.NewAppStoreKeyGeneralFromDataPiece(d)
	s.L.Lock()
	defer s.L.Unlock()

	if cm, ok := s.C[askg]; ok {
		started := cm.Started

		if (cm.Blocked && cm.Cur > 1) ||
			(cm.Blocked && cm.Cur == 1 && d.B() == repo.False && d.E() != repo.False) ||

			!cm.Blocked {
			//logger.L.Infof("in store.Dec cm before %v\n", cm)
			if d.IsSub() {
				cm.Cur--
				cm.Max--
				//logger.L.Infof("in store.Dec cm after decrementing %v\n", cm)
			} else {
				cm.Cur--

				if !cm.Started {
					cm.Started = true
				}
			}

			//logger.L.Infof("in store.Dec dataPiece with header %v made cm %v\n", d.GetHeader(), cm)
			delete(s.C, askg)
			s.C[askg] = cm
			//logger.L.Infof("in store.Dec s.C after %v\n", s.C)
			if cm.Cur == 0 && cm.Max-cm.Cur == 1 {
				//logger.L.Errorln("in store.Dec Cur == 0, resetting Store")

				if !started && cm.Started {

					if d.B() == repo.False && d.E() == repo.False {
						s.Reset(askg)
						return repo.FirstAndLast, nil
					}

					return repo.First, nil
				}
				if d.B() == repo.True {
					return repo.Intermediate, nil
				}
				// d.E() == repo.False

			}
			if !started && cm.Started { // may be different d.Part
				return repo.First, nil
			}
			if cm.Cur == cm.Max {
				return repo.Unordered, nil
			}
			if cm.Cur == 0 {

				if s.C[askg].Blocked {
					return repo.Intermediate, nil
				}
				// !Blocked
				s.Reset(askg)
				return repo.Last, nil
			}
			return repo.Intermediate, nil
		}
		return repo.Unordered, fmt.Errorf("in store.Dec cannot dec further")
	}
	return repo.Unordered, fmt.Errorf("in store.Dec askg \"%v\" not found", askg)
}
func (s *StoreStruct) Counter(askg repo.AppStoreKeyGeneral) int {
	s.L.RLock()
	defer s.L.RUnlock()
	return s.C[askg].Cur
}

func (s *StoreStruct) Reset(askg repo.AppStoreKeyGeneral) {
	//logger.L.Infof("store.Reset invoked with askg: %v, s.R: %v, s.B: %v, s.C: %v\n", askg, s.R, s.B, s.C)

	if _, ok := s.R[askg]; ok {
		if len(s.R) < 2 {
			s.R = make(map[repo.AppStoreKeyGeneral]map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue)
		}
		delete(s.R, askg)
	}

	if _, ok := s.B[askg]; ok {
		if len(s.B) < 2 {
			s.B = make(map[repo.AppStoreKeyGeneral][]repo.DataPiece)
		}
		delete(s.B, askg)
	}

	if _, ok := s.C[askg]; ok {
		//logger.L.Infof("in store.Reset s.c: %v, len: %d\n", s.C, len(s.C))
		if len(s.C) < 2 {

			s.C = make(map[repo.AppStoreKeyGeneral]repo.Counter)
			//logger.L.Infof("in store.Reset after deleting s.C: %v\n", s.C)
		}
		delete(s.C, askg)
	}
	/*
	   c := 0
	   	for i := range s.C {
	   		if s.C[i].Cur == 0 {
	   			c++
	   		}
	   	}
	   	if c == len(s.C) {
	   		s.C = map[repo.AppStoreKeyGeneral]repo.Counter{}
	   	}
	*/
}

func Equal(a, b repo.DataPiece) bool {
	if a.TS() == b.TS() &&
		a.Part() == b.Part() &&
		a.B() == b.B() &&
		a.E() == b.E() &&
		bytes.Contains(a.GetBody(0), b.GetBody(0)) {
		return true
	}
	return false
}

func isIn(i int, d []int) bool {
	for _, v := range d {
		if i == v {
			return true
		}
	}
	return false
}

func (s *StoreStruct) Presence(d repo.DataPiece) (repo.Presense, error) {
	askg, askd, vv := repo.NewAppStoreKeyGeneralFromDataPiece(d), repo.NewAppStoreKeyDetailed(d), make(map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue)
	//logger.L.Infof("store.Presense invoked by dataPiece's header %v, s.R: %v, s.C[askg]: %v\n", d.GetHeader(), s.R, s.C[askg])
	if m1, ok := s.R[askg]; ok {
		//logger.L.Infoln("store.Presense askg met")
		//logger.L.Infof("store.Presense askd.T(): %v, m1: %v\n", askd.T(), m1)
		if m2, ok := m1[askd]; ok && d.B() == repo.True {
			//logger.L.Infoln("store.Presense askd met")
			if s.C[askg].Cur == 1 && s.C[askg].Blocked {
				return repo.Presense{}, fmt.Errorf("in store.Presense matched but Cur == 1 && Blocked")
			}
			vv[askd.F()] = m2
			//logger.L.Infof("store.Presense vv: %v\n", vv)
			if m2t, ok := m1[askd.T()]; ok && d.E() == repo.Probably {
				//logger.L.Infoln("store.Presense d.E() == repo.Probably, true branch met")
				//logger.L.Infof("store.Presense askd.T(): %v, m1: %v\n", askd.T(), m1)
				vv[askd.T()] = m2t
				//logger.L.Infof("store.Presense vv: %v\n", vv)
				return repo.NewPresense(true, true, true, vv), nil
			}
			return repo.NewPresense(true, true, false, vv), nil
		}
		if d.IsSub() {
			//logger.L.Infoln("store.Presense AppSub branch")
			if m2f, ok := s.R[askg][askd.F()]; ok && m2f[false].E == repo.Probably {
				//logger.L.Infof("store.Presense AppSub branch m2f: %v", m2f)
				vv[askd.F()] = m2f
				return repo.NewPresense(true, true, true, vv), nil
			}
			//logger.L.Infof("in store.Presense AppSub case s.R: %v\n", s.R)
			return repo.NewPresense(true, false, false, nil), nil
		}
		if d.B() == repo.False && d.E() == repo.Probably {
			if m2t, ok := s.R[askg][askd.T()]; ok && m2t[true].E == repo.Probably {
				//logger.L.Infof("store.Presense AppSub branch m2f: %v", m2f)
				vv[askd.T()] = m2t
				return repo.NewPresense(true, true, true, vv), nil
			}
			return repo.NewPresense(true, true, false, nil), nil
		}
		return repo.NewPresense(true, false, false, nil), nil
	}
	if d.B() == repo.False && d.E() == repo.Probably {

	}
	return repo.Presense{}, nil
}

func (s *StoreStruct) Act(d repo.DataPiece, sc repo.StoreChange) {
	//logger.L.Infof("store.ACT is invoked for dataPiece header %v, sc.A: %d, s.R: %v, s.B: %v\n", d.GetHeader(), sc.A, s.R, s.B)
	//logger.L.Infof("store.ACT is invoked for dataPiece with body %q, s.R: %v, s.B: %v\n", d.GetBody(0), s.R, s.B)
	askg, vv := repo.NewAppStoreKeyGeneralFromDataPiece(d), make(map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue)

	switch sc.A {
	case repo.Change:
		if m1, ok := s.R[askg]; ok { // ASKD met
			//logger.L.Infof("in store.Act repo.NewAppStoreKeyDetailed(d): %v\n", repo.NewAppStoreKeyDetailed(d))

			for i, _ := range sc.From {
				if _, ok := m1[i]; ok {
					//logger.L.Infof("in store.Act deleting v: %v\n", v)
					delete(s.R[askg], i)
				}
				//logger.L.Infof("in store.Act after deleting s.R: %v\n", s.R)

			}
			for i, v := range sc.To {
				s.R[askg][i] = v
				//logger.L.Infof("in store.Act dataPiece header %v ==> s.R: %v\n", d.GetHeader(), s.R)
				return
			}
		}
		// no ASKG
		if d.B() == repo.False {
			for i, v := range sc.To {
				vv[i] = v
				s.R[askg] = vv
				//logger.L.Infof("in store.Act after adding s.R: %v\n", s.R)
				//logger.L.Infof("in store.Act dataPiece header %v ==> s.R: %v\n", d.GetHeader(), s.R)
				return
			}
		}
		s.R[askg] = make(map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue)

	case repo.Buffer:
		s.BufferAdd(d)
	case repo.Remove:
		delete(s.R, askg)
	}
}

func (s *StoreStruct) Update(dr repo.DetailedRecord) {

	if len(dr.DR) == 0 || len(dr.DR) > 2 {
		return
	}
	askg := repo.AppStoreKeyGeneral{TS: dr.ASKD.SK.TS}

	if _, ok := s.R[askg]; !ok {
		vv := make(map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue, 0)
		vv[dr.ASKD.IncPart()] = dr.DR
		s.R[askg] = vv
		return
	}
	delete(s.R[askg], dr.ASKD)
	s.R[askg][dr.ASKD.IncPart()] = dr.DR
}

func (s *StoreStruct) Delete(askg repo.AppStoreKeyGeneral) {
	delete(s.R, askg)
}

func (s *StoreStruct) Unblock(askg repo.AppStoreKeyGeneral) {
	s.L.Lock()
	defer s.L.Unlock()

	mr := s.C[askg]
	mr.Blocked = s.L.TryLock()

	delete(s.C, askg)

	s.C[askg] = mr
}
