package utils

import (
	"reflect"
)

// compare two arrays. find Missing elements in one or the other
type ArrayComparer struct {
	elements_in_array_1_but_not_2 []int
	elements_in_array_2_but_not_1 []int
}

// TODO: check if we can use reflect.ValueOf(array1).Index(i) to remove the need for external comp function
func CompareArray(array1, array2 interface{}, comp func(i, j int) bool) *ArrayComparer {
	res := &ArrayComparer{}
	l1 := reflect.ValueOf(array1).Len()
	l2 := reflect.ValueOf(array2).Len()
	// check for elements in array1 that are not in array2
	for i := 0; i < l1; i++ {
		// check if element 'i' from 'array1' is in 'array2'
		found := false
		for j := 0; j < l2; j++ {
			if comp(i, j) {
				found = true
				break
			}
		}
		if !found {
			res.elements_in_array_1_but_not_2 = append(res.elements_in_array_1_but_not_2, i)
		}
	}

	// check for elements in array2 that are not in array1
	for j := 0; j < l2; j++ {
		// check if element 'j' from 'array2' is in 'array1'
		found := false
		for i := 0; i < l1; i++ {
			if comp(i, j) {
				found = true
				break
			}
		}
		if !found {
			res.elements_in_array_2_but_not_1 = append(res.elements_in_array_2_but_not_1, j)
		}
	}

	return res
}

// returns true if both arrays contain the same elements (disregarding the order they are in)
func (ac *ArrayComparer) EqualElements() bool {
	if len(ac.elements_in_array_1_but_not_2) == 0 && len(ac.elements_in_array_2_but_not_1) == 0 {
		return true
	}
	return false
}

// returns indices of elements that are in array 1 but not array 2
func (ac *ArrayComparer) ElementsIn1ButNot2() []int {
	return ac.elements_in_array_1_but_not_2
}

// returns indices of elements that are in array 2 but not array 1
func (ac *ArrayComparer) ElementsIn2ButNot1() []int {
	return ac.elements_in_array_2_but_not_1
}
