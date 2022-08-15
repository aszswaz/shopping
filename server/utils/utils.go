package utils

type StringSlice struct {
	Slice []string
}

func (ss *StringSlice) Contain(element string) bool {
	if ss.Slice == nil {
		return false
	}
	for _, item := range ss.Slice {
		if item == element {
			return true
		}
	}
	return false
}

func (ss *StringSlice) Deduplication() *StringSlice {
	result := StringSlice{make([]string, 0)}
	for _, item := range ss.Slice {
		if !result.Contain(item) {
			result.Slice = append(result.Slice, item)
		}
	}
	ss.Slice = result.Slice
	return ss
}
