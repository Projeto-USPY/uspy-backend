package emulator

import "testing"

func GetEmulator(t *testing.T) {
	if _, err := Get(); err != nil {
		t.Fatal(err)
	}
}

func TestEmulator(t *testing.T) {
	t.Run("Get emulator", GetEmulator)
	t.Run("Get emulator", GetEmulator)
	t.Run("Get emulator", GetEmulator)

	if insertCnt > 1 {
		t.Fatal("InsertUtilsOnce ran more than once")
	}
}
