package dataHandler

import "github.com/vynovikov/highLoadParser/internal/logger"

type memoryDataHandlerStruct struct {
	Map    map[key]map[bool]value
	Buffer []DataPiece
}

func NewMemoryDataHandler() *memoryDataHandlerStruct {
	return &memoryDataHandlerStruct{}
}

func (m *memoryDataHandlerStruct) Create(d DataPiece) error {
	logger.L.Printf("in dataHandler creating dataPiece = %v\n", d)
	return nil
}

func (m *memoryDataHandlerStruct) Read(DataPiece) (value, error) {
	return value{}, nil
}

func (m *memoryDataHandlerStruct) Updade(DataPiece) error {
	return nil
}

func (m *memoryDataHandlerStruct) Delete(DataPiece) error {
	return nil
}

func (m *memoryDataHandlerStruct) Check(d DataPiece) (Presence, error) {

	//mapKey, mapVal := key{}, value{}

	logger.L.Infof("in dataHandler.Check d: %#v", d)

	return Presence{}, nil
}

/*
func (s *StoreStruct) Presence(d repo.DataPiece) (repo.Presense, error) {
	askg, askd, vv := repo.NewAppStoreKeyGeneralFromDataPiece(d), repo.NewAppStoreKeyDetailed(d), make(map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue)
	if m1, ok := s.R[askg]; ok {
		if m2, ok := m1[askd]; ok && d.B() == repo.True {
			if s.C[askg].Cur == 1 && s.C[askg].Blocked {
				return repo.Presense{}, fmt.Errorf("in store.Presense matched but Cur == 1 && Blocked")
			}
			vv[askd.F()] = m2
			if m2t, ok := m1[askd.T()]; ok && d.E() == repo.Probably {
				vv[askd.T()] = m2t
				return repo.NewPresense(true, true, true, vv), nil
			}
			return repo.NewPresense(true, true, false, vv), nil
		}
		if d.IsSub() {
			if m2f, ok := s.R[askg][askd.F()]; ok && m2f[false].E == repo.Probably {
				vv[askd.F()] = m2f
				return repo.NewPresense(true, true, true, vv), nil
			}
			return repo.NewPresense(true, false, false, nil), nil
		}
		if d.B() == repo.False && d.E() == repo.Probably {
			if m2t, ok := s.R[askg][askd.T()]; ok && m2t[true].E == repo.Probably {
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
*/
