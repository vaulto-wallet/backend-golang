package helpers

import (
	"fmt"
	"strings"
)

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

func UintFind(array []uint, value uint) int {
	for i, v := range array {
		if v == value {
			return i
		}
	}
	return -1
}

func UintToString(array []uint) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(array)), ","), "[]")
}

func Remove(array []uint, value uint) (ret []uint) {
	copy(ret, array)
	idx := UintFind(array, value)
	if idx == -1 {
		return
	}
	copy(ret[idx:], ret[idx+1:])
	ret = ret[:len(ret)-1]
	return
}

func UintAppendNewArray(array []uint, value []uint) []uint {
	for _, a := range value {
		if !UintInArray(array, a) {
			array = append(array, a)
		}
	}
	return array
}
