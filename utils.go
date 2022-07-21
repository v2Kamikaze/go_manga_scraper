package main

import "reflect"

func Contains(v []Chapter, e Chapter) bool {
	for i := range v {
		if reflect.DeepEqual(e, v[i]) {
			return true
		}
	}
	return false
}

func IndexOf(v []any, e any) int {
	for i := range v {
		if reflect.DeepEqual(e, v[i]) {
			return i
		}
	}
	return -1
}
