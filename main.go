package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

func readFilesFromDir(dirPath string) (map[string]string, error) {
	files, err := filepath.Glob(dirPath + "/*.txt")
	if err != nil {
		return nil, err
	}
	texts := make(map[string]string)
	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatalf("Error reading file %s: %v", file, err)
		}
		texts[file] = string(content)
	}
	return texts, nil
}

// calcula la distancia Levenshtein
func levenshteinDistance(s1, s2 string) int {
	r1, r2 := utf8.RuneCountInString(s1), utf8.RuneCountInString(s2)
	dp := make([][]int, r1+1)
	for i := range dp {
		dp[i] = make([]int, r2+1)
		dp[i][0] = i
	}
	for j := range dp[0] {
		dp[0][j] = j
	}

	for i, ri := 1, []rune(s1); i <= r1; i++ {
		for j, rj := 1, []rune(s2); j <= r2; j++ {
			cost := 0
			if ri[i-1] != rj[j-1] {
				cost = 1
			}
			dp[i][j] = min(min(dp[i-1][j]+1, dp[i][j-1]+1), dp[i-1][j-1]+cost)
		}
	}
	return dp[r1][r2]
}

// Encuentra todas las subcadenas comunes
func findAllCommonSubstrings(text1, text2 string) []string {
	substrings := make(map[string]struct{})
	n := len(text1)
	m := len(text2)

	// Algoritmo de Fuerza Bruta para encontrar todas las subcadenas comunes
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if text1[i] == text2[j] {
				length := 0
				for i+length < n && j+length < m && text1[i+length] == text2[j+length] {
					substring := text1[i : i+length+1]
					substrings[substring] = struct{}{}
					length++
				}
			}
		}
	}

	// Convertir map a slice
	result := make([]string, 0, len(substrings))
	for substring := range substrings {
		result = append(result, substring)
	}
	return result
}

func highlightCommonSubstrings(text string, substrings []string) string {
	highlighted := text
	for _, substring := range substrings {
		highlighted = strings.ReplaceAll(highlighted, substring, `<span style="background-color: yellow; color: #d3163b;">`+substring+`</span>`)
	}
	return highlighted
}

// calcula la similitud usando todas las subcadenas comunes y la distancia de Levenshtein
func calculateSimilarity(text1, text2 string) float64 {
	// Encuentra todas las subcadenas comunes
	commonSubstrings := findAllCommonSubstrings(text1, text2)
	totalLength := 0
	for _, substring := range commonSubstrings {
		totalLength += len(substring)
	}
	longestLength := max(len(text1), len(text2))
	substringsSimilarity := float64(totalLength) / float64(longestLength)

	// Calcula la distancia Levenshtein
	editDist := levenshteinDistance(text1, text2)
	normalizedEditDist := 1.0 - float64(editDist)/float64(longestLength) // Normalize edit distance

	// Combina ambas similitudes (50% cada una)
	return 0.5*substringsSimilarity + 0.5*normalizedEditDist
}

func mergeSortPairs(pairs []TextPair) []TextPair {
	if len(pairs) <= 1 {
		return pairs
	}
	mid := len(pairs) / 2
	left := mergeSortPairs(pairs[:mid])
	right := mergeSortPairs(pairs[mid:])
	return mergePairs(left, right)
}

func mergePairs(left, right []TextPair) []TextPair {
	result := []TextPair{}
	i, j := 0, 0
	for i < len(left) && j < len(right) {
		if left[i].Similarity >= right[j].Similarity {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}
	result = append(result, left[i:]...)
	result = append(result, right[j:]...)
	return result
}

type TextPair struct {
	File1       string
	File2       string
	Similarity  float64
	Highlighted string
}

// Frontend
func generateHTML(pairs []TextPair) {
	html := `<html>
<head>
<title>Plagiarism Detector</title>
<style>
/* Aquí va el CSS que ya tenías */
</style>
</head>
<body>
<h1>Plagiarism Detector •*.✸ </h1>
<h3 style="text-align: center;">Web App in Golang that compares the content of different texts about the same topic, calculating and highlighting their similarity using multiple matching substrings.</h3>`

	counter := 1
	for _, pair := range pairs {
		splitIndex := strings.Index(pair.Highlighted, "</pre><pre>") + 6
		firstText := pair.Highlighted[:splitIndex]
		secondText := pair.Highlighted[splitIndex:]

		html += fmt.Sprintf(`
        <div class="container">
            <description>%d. %s vs %s</description>
            <h2>Similarity: %.2f%%</h2>
            <div class="text-container">
                <div class="text-content">
                    <pre>%s</pre>
                </div>
                <div class="text-content">
                    <pre>%s</pre>
                </div>
            </div>
        </div>`,
			counter, filepath.Base(pair.File1), filepath.Base(pair.File2), pair.Similarity*100, firstText, secondText)
		counter++
	}
	html += "</body></html>"
	err := ioutil.WriteFile("plagiarism_report.html", []byte(html), 0644)
	if err != nil {
		log.Fatalf("Error writing HTML file: %v", err)
	}
	fmt.Println("HTML report generated: plagiarism_report.html")
}

// subraya
func calculateSimilarityAndHighlight(texts map[string]string) {
	keys := make([]string, 0, len(texts))
	for key := range texts {
		keys = append(keys, key)
	}
	var pairs []TextPair
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			text1 := texts[keys[i]]
			text2 := texts[keys[j]]
			similarity := calculateSimilarity(text1, text2)
			commonSubstrings := findAllCommonSubstrings(text1, text2)
			highlightedText1 := highlightCommonSubstrings(text1, commonSubstrings)
			highlightedText2 := highlightCommonSubstrings(text2, commonSubstrings)
			pairs = append(pairs, TextPair{
				File1:       keys[i],
				File2:       keys[j],
				Similarity:  similarity,
				Highlighted: fmt.Sprintf("<pre>%s</pre><pre>%s</pre>", highlightedText1, highlightedText2),
			})
		}
	}
	sortedPairs := mergeSortPairs(pairs)
	if len(sortedPairs) > 10 {
		sortedPairs = sortedPairs[:10]
	}
	generateHTML(sortedPairs)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	texts, err := readFilesFromDir("dataset")
	if err != nil {
		log.Fatal(err)
	}

	// preview en la terminal
	for fileName, content := range texts {
		fmt.Printf("File: %s\nContent: %.50s...\n", fileName, content)
	}

	// hace el HTML
	calculateSimilarityAndHighlight(texts)
}