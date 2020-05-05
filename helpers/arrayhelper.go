package helpers

func UintInArray(array []uint, value uint) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}

func UintAppendNew(array []uint, value uint) []uint {
	if !UintInArray(array, value) {
		array = append(array, value)
	}
	return array
}
