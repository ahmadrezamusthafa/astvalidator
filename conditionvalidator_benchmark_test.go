package astvalidator

import "testing"

//BENCHMARK ValidateCondition
//Improvement history:
//------------------------------------
//	attempt	   |  time per loop
//------------------------------------
//  4878042	      243 ns/op (now)
//------------------------------------
func BenchmarkValidateCondition(b *testing.B) {
	referenceQuery := `id=1 && ( division = engineering || division = finance )`
	input := `id=1 && division = engineering`
	referenceCondition, _ := GenerateCondition(referenceQuery)
	inputCondition, _ := GenerateCondition(input)

	for n := 0; n < b.N; n++ {
		referenceCondition.ValidateCondition(inputCondition)
	}
}
