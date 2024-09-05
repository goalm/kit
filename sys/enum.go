package sys

// bi-directional map
type Enum struct {
	name map[int]string // idx -> name
	idx  map[string]int // name -> idx
	desc map[int]string // idx -> desc
}

func NewEnum() *Enum {
	return &Enum{
		name: make(map[int]string),
		idx:  make(map[string]int),
		desc: make(map[int]string),
	}
}

func (m *Enum) Add(idx int, name, desc string) {
	// optionally verify uniqueness constraint
	m.name[idx] = name
	m.idx[name] = idx
	m.desc[idx] = desc
}

func (m *Enum) IntToStr(idx int) (string, bool) {
	b, ok := m.name[idx]
	return b, ok
}

func (m *Enum) StrToInt(name string) (int, bool) {
	a, ok := m.idx[name]
	return a, ok
}

func (m *Enum) Size() int {
	return len(m.name)
}
