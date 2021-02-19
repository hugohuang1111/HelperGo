package utils

import "strings"

func RemoveArrayString(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func InsertInterface2Array(a []interface{}, index int, value interface{}) []interface{} {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	return a
}

func InsertString2Array(a []string, index int, value string) []string {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	return a
}

func StringContainArray(s string, arr []string) int {
	for idx, str := range arr {
		if strings.Contains(s, str) {
			return idx
		}
	}

	return -1
}

func ArrayContainString(s string, arr []string) int {
	for idx, str := range arr {
		if strings.Contains(str, s) {
			return idx
		}
	}

	return -1
}
