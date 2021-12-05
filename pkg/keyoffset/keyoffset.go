package keyoffset

func getInterval(index []KeyOffset, key string) (*KeyOffset, *KeyOffset) {
	if len(index) == 0 {
		return nil, nil
	}
	
	if len(index) == 1 {
		if key < index[0].Key {
			return nil, &index[0]
		} else {
			return &index[0], nil
		}
	}

	mid := len(index) / 2

	var l, r *KeyOffset
	if key < index[mid].Key {
		l, r = getInterval(index[:mid], key)
		if r == nil {
			r = &index[mid]
		}
	} else {
		l, r = getInterval(index[mid:], key)
		if l == nil && mid > 0 {
			l = &index[mid-1]
		}
	}

	return l, r
}

type KeyOffset struct {
	Key    string
	Offset int64
}