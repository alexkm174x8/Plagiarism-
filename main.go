package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"
)

type Suffix struct {
	Index int
	Text  string
}

func BuildSuffixArray(text string) []int {
	n := len(text)
	suffixes := make([]Suffix, n)

	for i := 0; i < n; i++ {
		suffixes[i] = Suffix{Index: i, Text: text[i:]}
	}

	sort.Slice(suffixes, func(i, j int) bool {
		return suffixes[i].Text < suffixes[j].Text
	})

	suffixArray := make([]int, n)
	for i := 0; i < n; i++ {
		suffixArray[i] = suffixes[i].Index
	}

	return suffixArray
}

func BuildLCPArray(text string, suffixArray []int) []int {
	n := len(text)
	rank := make([]int, n)
	lcp := make([]int, n)

	for i := 0; i < n; i++ {
		rank[suffixArray[i]] = i
	}

	h := 0
	for i := 0; i < n; i++ {
		if rank[i] > 0 {
			j := suffixArray[rank[i]-1]
			for i+h < n && j+h < n && text[i+h] == text[j+h] {
				h++
			}
			lcp[rank[i]] = h
			if h > 0 {
				h--
			}
		}
	}

	return lcp
}

func CompareDocuments(text1, text2 string) float64 {
	combinedText := text1 + "#" + text2
	suffixArray := BuildSuffixArray(combinedText)
	lcpArray := BuildLCPArray(combinedText, suffixArray)

	maxLCP := 0
	for _, lcp := range lcpArray {
		if lcp > maxLCP {
			maxLCP = lcp
		}
	}

	minLength := len(text1)
	if len(text2) < minLength {
		minLength = len(text2)
	}
	similarity := float64(maxLCP) / float64(minLength) * 100.0

	return similarity
}

func ReadFile(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func main() {
	doc1, err := ReadFile("dataset/g0pA_taska.txt")
	if err != nil {
		log.Fatalf("Error al leer el archivo doc1.txt: %v", err)
	}

	doc2, err := ReadFile("dataset/g0pA_taskb.txt")
	if err != nil {
		log.Fatalf("Error al leer el archivo doc2.txt: %v", err)
	}

	similarity := CompareDocuments(doc1, doc2)
	fmt.Printf("La similitud entre los documentos es: %.2f%%\n", similarity)
}
