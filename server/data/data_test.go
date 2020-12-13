package data

import (
	"testing"
)

func TestLoadData(t *testing.T) {
	LoadData()
	t.Log(courses, professors)
}
