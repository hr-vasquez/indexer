package main

import (
	"testing"
)

func TestIndexer(t *testing.T) {
	// Arrange
	fileToIndex := "./sample"

	// Act
	res := IndexSource(fileToIndex)

	// Assert
	if res.StatusCode != 200 {
		t.Fatalf("expected %d, but got %d.", 200, res.StatusCode)
	}
}

func BenchmarkIndexer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IndexSource("./sample")
	}
}
