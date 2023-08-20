package map_provider

type Variable struct {
	Name  string
	Value MapProvider
}

func (v *Variable) Find(name string) (string, bool) {
	return v.Value.Find(name)
}

var _ MapProvider = &Variable{}
