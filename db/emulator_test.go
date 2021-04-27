package db

import "testing"

func TestGetEmulator(t *testing.T) {
	if _, err := GetEmulator(); err != nil {
		t.Fatal(err)
	}
}
