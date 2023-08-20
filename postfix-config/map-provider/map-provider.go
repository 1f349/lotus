package map_provider

// MapProvider is an interface to allow looking up mapped values from variables,
// hash files or mysql queries.
type MapProvider interface {
	Find(name string) (string, bool)
}

// SequenceMapProvider calls Find against each provider in a slice and outputs
// the true mapped value of the input. first mapped value found. If the input was
// not found then "", false is returned.
type SequenceMapProvider []MapProvider

func (s SequenceMapProvider) Find(name string) (string, bool) {
	for _, i := range s {
		if find, ok := i.Find(name); ok {
			return find, true
		}
	}
	return "", false
}

var _ MapProvider = SequenceMapProvider{}
