package flaeg

import (
	"reflect"
	"testing"
)

func TestSliceStringsSet(t *testing.T) {
	checkMap := map[string]SliceStrings{
		"str":            SliceStrings{"str"},
		"str1,str2":      SliceStrings{"str1", "str2"},
		"str1;str2":      SliceStrings{"str1", "str2"},
		"str1,str2;str3": SliceStrings{"str1", "str2", "str3"},
	}
	for str, check := range checkMap {
		var slice SliceStrings
		if err := slice.Set(str); err != nil {
			t.Fatalf("Error :%s", err)
		}
		if !reflect.DeepEqual(slice, check) {
			t.Fatalf("Expected:%s\ngot:%s", check, slice)
		}
	}
}
func TestSliceStringsSetAdd(t *testing.T) {
	slice := SliceStrings{"str1"}
	//test
	if err := slice.Set("str2,str3"); err != nil {
		t.Fatalf("Error :%s", err)
	}
	//check
	check := SliceStrings{"str1", "str2", "str3"}
	if !reflect.DeepEqual(slice, check) {
		t.Fatalf("Expected:%s\ngot:%s", check, slice)
	}
}

func TestSliceStringsGet(t *testing.T) {
	slices := []SliceStrings{
		SliceStrings{"str"},
		SliceStrings{"str1", "str2"},
		SliceStrings{"str1", "str2", "str3"},
	}
	check := [][]string{
		[]string{"str"},
		[]string{"str1", "str2"},
		[]string{"str1", "str2", "str3"},
	}
	for i, slice := range slices {
		if !reflect.DeepEqual(slice.Get(), check[i]) {
			t.Fatalf("Expected:%s\ngot:%s", check[i], slice)
		}
	}
}

func TestSliceStringsString(t *testing.T) {
	slices := []SliceStrings{
		SliceStrings{"str"},
		SliceStrings{"str1", "str2"},
		SliceStrings{"str1", "str2", "str3"},
	}
	check := []string{
		"[str]",
		"[str1 str2]",
		"[str1 str2 str3]",
	}
	for i, slice := range slices {
		if !reflect.DeepEqual(slice.String(), check[i]) {
			t.Fatalf("Expected:%s\ngot:%s", check[i], slice)
		}
	}
}

func TestSliceStringsSetValue(t *testing.T) {
	check := []SliceStrings{
		SliceStrings{"str"},
		SliceStrings{"str1", "str2"},
		SliceStrings{"str1", "str2", "str3"},
	}
	slices := [][]string{
		[]string{"str"},
		[]string{"str1", "str2"},
		[]string{"str1", "str2", "str3"},
	}
	for i, s := range slices {
		var slice SliceStrings
		slice.SetValue(s)
		if !reflect.DeepEqual(slice, check[i]) {
			t.Fatalf("Expected:%s\ngot:%s", check[i], slice)
		}
	}
}
