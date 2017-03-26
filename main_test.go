package main

import "testing"

func BenchmarkDict1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		process("dict1", "out1")
	}
}

func BenchmarkDict2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		process("dict2", "out2")
	}
}
