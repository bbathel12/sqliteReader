package main

import "testing"

func BenchmarkWriteToFile(b *testing.B) {

	for i := 0; i < b.N; i++ {
		testChan := make(chan string, 1)
		testChan <- "96d8b2f9e0f2e8c1aea0f998cf8f9224"
		close(testChan)
		writeToFile("./output.txt", testChan)
	}

}
