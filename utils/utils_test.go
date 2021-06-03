/* Package db contains useful functions related to the Firestore Database */
package utils

import "testing"

func BenchmarkHashPasswords(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hash, _ := Bcrypt("SenhaU3l34178!Fodida18723@#!")
		b.Log(hash)
	}
}
