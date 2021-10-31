package utils

import (
	"reflect"

	"cloud.google.com/go/firestore"
)

// MergeWithout takes an object and a set of firestore fields to return a SetOption without them
//
// It is like firestore.MergeAll, but it allows some fields to be excluded
func MergeWithout(value interface{}, fields ...string) firestore.SetOption {
	f := reflect.TypeOf(value)

	// make set from array of fields
	set := make(map[string]struct{})
	for _, v := range fields {
		set[v] = struct{}{}
	}

	results := make([]string, 0)
	for i := 0; i < f.NumField(); i++ {
		tag := f.Field(i).Tag.Get("firestore")

		if tag == "-" { // if field should be ignored
			continue
		}

		if _, ok := set[tag]; ok { // if tag should excluded
			continue
		}

		results = append(results, tag)
	}

	return firestore.Merge(results)
}
