package dataHandler

type memoryDataHandlerStruct struct {
	Map    map[key]map[bool]value
	Buffer []DataHandlerDTO
}

func NewMemoryDataHandler() *memoryDataHandlerStruct {
	return &memoryDataHandlerStruct{
		Map:    make(map[key]map[bool]value),
		Buffer: make([]DataHandlerDTO, 0),
	}
}

func (m *memoryDataHandlerStruct) Create(d DataHandlerDTO) error {

	key := newKey(d)
	val := newValue(d)

	switch len(m.Map[key]) {

	case 0: // Only possible if !d.IsSub

		key.Part++

		valMap := make(map[bool]value)

		valMap[false] = val

		m.Map[key] = valMap

	default:
	}

	return nil
}

func (m *memoryDataHandlerStruct) Read(DataHandlerDTO) (value, error) {
	return value{}, nil
}

func (m *memoryDataHandlerStruct) Updade(DataHandlerDTO) error {
	return nil
}

func (m *memoryDataHandlerStruct) Delete(string) error {
	return nil
}

func (m *memoryDataHandlerStruct) Check(d DataHandlerDTO) (Presence, error) {

	key := newKey(d)

	if v, ok := m.Map[key]; ok {
		return Presence{value: v}, nil
	}

	return Presence{}, nil
}

/*
func (s *StoreStruct) Presence(d repo.DataHandlerDTO) (repo.Presense, error) {
	askg, askd, vv := repo.NewAppStoreKeyGeneralFromDataHandlerDTO(d), repo.NewAppStoreKeyDetailed(d), make(map[repo.AppStoreKeyDetailed]map[bool]repo.AppStoreValue)
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
