package util

import (
	"testing"
)

/*
	Feature #A:
		The string generator generates only hex values. Providing Unit tests and benchmarks for the renamed 'RandHexString' function.
		Unit tests can be executed from this directory with the following command:
			go test -v
		Benchmarks can be executed from this directory with the following command:
			go test -bench=.
*/

// Testing if the length of the generated string is right
func TestRandHexStringLength(t *testing.T) {
	lengths_of_strings := []int{0, 5, 10, 55, 84}
	for _, length := range lengths_of_strings {
		result := RandHexString(length)
		if len(result) != length {
			t.Errorf("Expected length %d, but got %d", length, len(result))
		}
	}
}

// Testing if the generated string contains only hexadecimal characters
func TestRandHexStringContent(t *testing.T) {
	length_of_string := 100
	result := RandHexString(length_of_string)
	for _, c := range result {
		if c < '0' || (c > '9' && c < 'a') || c > 'f' {
			t.Errorf("Expected only hexadecimal characters, but got %c", c)
		}
	}
}

// Testing if the generated string is deterministic for the same seed
func TestRandHexStringDeterministic(t *testing.T) {
	length_of_string := 10
	randx.Seed(42)
	result1 := RandHexString(length_of_string)
	randx.Seed(42)
	result2 := RandHexString(length_of_string)
	if result1 != result2 {
		t.Errorf("Expected deterministic result, but got different results: %s and %s", result1, result2)
	}
}

// Benchmarks the RandHexString function
func BenchmarkRandHexString(b *testing.B) {
	length_of_string := 10
	for i := 0; i < b.N; i++ {
		RandHexString(length_of_string)
	}
}
