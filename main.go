package main

import (
	"fmt"
	"sort"
)

type SuffixArrayEntry struct {
	index  int
	suffix string
}

func ConstruirSuffixArray(texto string) []int {
	n := len(texto)
	suffixes := make([]SuffixArrayEntry, n)

	for i := 0; i < n; i++ {
		suffixes[i] = SuffixArrayEntry{i, texto[i:]}
	}

	sort.Slice(suffixes, func(i, j int) bool {
		return suffixes[i].suffix < suffixes[j].suffix
	})

	suffixArray := make([]int, n)
	for i := 0; i < n; i++ {
		suffixArray[i] = suffixes[i].index
	}

	return suffixArray
}

func LCPArray(texto string, suffixArray []int) []int {
	n := len(suffixArray)
	lcp := make([]int, n)
	invSuffix := make([]int, n)

	for i := 0; i < n; i++ {
		invSuffix[suffixArray[i]] = i
	}

	k := 0
	for i := 0; i < n; i++ {
		if invSuffix[i] == n-1 {
			k = 0
			continue
		}

		j := suffixArray[invSuffix[i]+1]
		for i+k < n && j+k < n && texto[i+k] == texto[j+k] {
			k++
		}

		lcp[invSuffix[i]] = k
		if k > 0 {
			k--
		}
	}

	return lcp
}

func CompararTextos(texto1, texto2 string) {
	fmt.Println("Texto 1:", texto1)
	fmt.Println("Texto 2:", texto2)

	texto := texto1 + "$" + texto2 + "#"
	fmt.Println("\nTexto concatenado:", texto)

	suffixArray := ConstruirSuffixArray(texto)
	fmt.Println("\nSuffix Array:", suffixArray)

	lcp := LCPArray(texto, suffixArray)
	fmt.Println("\nLCP Array:", lcp)
}

func main() {
	texto1 := "banana"
	texto2 := "bandana"
	CompararTextos(texto1, texto2)
}
