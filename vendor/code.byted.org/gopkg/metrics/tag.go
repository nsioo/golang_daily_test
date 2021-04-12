package metrics

import (
	"fmt"
)

type T struct {
	Name  string
	Value string
}

func Map2Tags(m map[string]string) []T {
	tt := make([]T, 0, len(m))
	for k, v := range m {
		tt = append(tt, Tag(k, v))
	}
	sortTags(tt) // go randomizes the iteration order, sort it
	return tt
}

func Tag(name, value string) T {
	return T{Name: name, Value: value}
}

func TFromI(name string, v interface{}) T {
	return T{Name: name, Value: fmt.Sprint(v)}
}

func (t T) less(o T) bool {
	if t.Name != o.Name {
		return t.Name < o.Name
	}
	return t.Value < o.Value
}

func sortTags(ss []T) {
	for i := 1; i < len(ss); i++ {
		for j := i; j > 0 && ss[j].less(ss[j-1]); j-- {
			ss[j], ss[j-1] = ss[j-1], ss[j]
		}
	}
}

func appendTags(b []byte, tt []T) []byte {
	sep := len(b) > 0
	for _, t := range tt {
		if !isValidString(t.Name) || !isValidString(t.Value) {
			continue
		}
		if sep {
			b = append(b, '|')
		}
		b = append(b, t.Name...)
		b = append(b, '=')
		b = append(b, t.Value...)
		sep = true
	}
	return b
}
