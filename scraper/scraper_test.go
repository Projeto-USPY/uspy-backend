package scraper

import (
	"fmt"
	"testing"
)

func TestScrapeDepartments(t *testing.T) {
	result := ScrapeDepartments()
	fmt.Println(result)
}

func TestScrapeCourseDescription(t *testing.T) {
	var course string = "SCC0200"
	desc, err := ScrapeCourseDescription(course)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(desc)
	}

	course = "SCC1234"
	desc, err = ScrapeCourseDescription(course)
	if err != nil {
		t.Log(err)
	} else {
		t.Fail()
	}
}
