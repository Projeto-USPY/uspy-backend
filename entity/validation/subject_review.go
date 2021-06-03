package validation

import (
	"reflect"
	"sort"

	"github.com/go-playground/validator/v10"
)

func getCategories() []string {
	return []string{"worth_it"}
}

func validateSubjectReview(f1 validator.FieldLevel) bool {
	keys := f1.Field().MapKeys()
	keysStr := make([]string, 0)

	for _, k := range keys {
		keysStr = append(keysStr, k.String())
	}

	categories := getCategories()

	if len(keys) != len(categories) {
		return false
	}

	sort.Strings(keysStr)
	sort.Strings(categories)

	for i := 0; i < len(keys); i++ {
		if !reflect.DeepEqual(keys[i].String(), categories[i]) {
			return false
		}
	}

	return true
}
