package main

import (
	"os"
	"testing"
)

func BenchmarkDay1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Day1()
	}
}

func BenchmarkDay2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Day2()
	}
}

func BenchmarkDay3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Day3()
	}
}

func BenchmarkDay4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Day4()
	}
}

func BenchmarkDay5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Day5()
	}
}

func TestMain(m *testing.M) {
	// Write to a file rather than io.Discard to be fair to a different
	// benchmark whose timings I'm comparing against.
	AocOut, _ = os.OpenFile("./trash", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	exitCode := m.Run()
	os.Exit(exitCode)
}
