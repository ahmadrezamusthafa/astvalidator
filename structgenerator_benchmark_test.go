package astvalidator

import "testing"

//BENCHMARK GenerateCondition
//Improvement history:
//------------------------------------
//	attempt	   |  time per loop
//------------------------------------
//  301020	      3341 ns/op
//  372768	      2793 ns/op (now)
//------------------------------------
func BenchmarkGenerateCondition(b *testing.B) {
	query := `id=1 && ( division = engineering || division = finance )`
	for n := 0; n < b.N; n++ {
		GenerateCondition(query)
	}
}

//BENCHMARK GetTokenAttributes
//Improvement history:
//------------------------------------
//	attempt	   |  time per loop
//------------------------------------
//  589063	      1962 ns/op
//  686304	      1792 ns/op (now)
//------------------------------------
func BenchmarkGetTokenAttributes(b *testing.B) {
	query := "id=1 && (division=engineering || division=finance)"
	for n := 0; n < b.N; n++ {
		getTokenAttributes(query)
	}
}
