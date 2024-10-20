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

// lcs
func longestCommonSubstring(text1, text2 string) (string, int) {
    combined := text1 + "#" + text2 + "$"
    suffixArray := buildSuffixArray(combined)
    n := len(combined)
    longest := ""

    for i := 0; i < n-1; i++ {
        if (suffixArray[i] < len(text1) && suffixArray[i+1] > len(text1)) ||
            (suffixArray[i] > len(text1) && suffixArray[i+1] < len(text1)) {
            lcs := commonPrefix(combined[suffixArray[i]:], combined[suffixArray[i+1]:])
            if len(lcs) > len(longest) {
                longest = lcs
            }
        }
    }
    return longest, len(longest)
}

func commonPrefix(s1, s2 string) string {
    n := len(s1)
    if len(s2) < n {
        n = len(s2)
    }
    for i := 0; i < n; i++ {
        if s1[i] != s2[i] {
            return s1[:i]
        }
    }
    return s1[:n]
}

//  suffix array 
func buildSuffixArray(text string) []int {
    n := len(text)
    suffixes := make([]string, n)
    suffixArray := make([]int, n)

    for i := 0; i < n; i++ {
        suffixes[i] = text[i:]
    }
    sortedSuffixes := mergeSortSuffixes(suffixes)
    for i := 0; i < n; i++ {
        suffixArray[i] = n - len(sortedSuffixes[i])
    }
    return suffixArray
}

// Merge sort para los suffixes
func mergeSortSuffixes(suffixes []string) []string {
    if len(suffixes) <= 1 {
        return suffixes
    }
    mid := len(suffixes) / 2
    left := mergeSortSuffixes(suffixes[:mid])
    right := mergeSortSuffixes(suffixes[mid:])
    return merge(left, right)
}

func merge(left, right []string) []string {
    result := []string{}
    i, j := 0, 0
    for i < len(left) && j < len(right) {
        if left[i] <= right[j] {
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

func highlightCommonSubstring(text, commonSubstring string) string {
    highlighted := strings.Replace(text, commonSubstring, `<span style="background-color: yellow; color: #d3163b;">`+commonSubstring+`</span>`, -1)
    return highlighted
}

// calcula la similitud usando lcs y lo de levisnon
func calculateSimilarity(text1, text2 string) float64 {
    //lcs
    _, lcsLength := longestCommonSubstring(text1, text2)
    longestLength := max(len(text1), len(text2))
    lcsSimilarity := float64(lcsLength) / float64(longestLength)

    // levensin
    editDist := levenshteinDistance(text1, text2)
    normalizedEditDist := 1.0 - float64(editDist)/float64(longestLength) // Normalize edit distance

    // las combina (50% cada una)
    return 0.5*lcsSimilarity + 0.5*normalizedEditDist
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
@font-face {
  font-family: Nohemi;
  src: url("Nohemi-Medium.otf")
}
@font-face {
  font-family: OffBit;
  src: url("OffBit.ttf")
}
body {
  background-color: #ffa1d7;
  font-family: Nohemi;
  color: #d3163b;
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}
.container {
  max-width: 1370px;
  margin: 20px auto;
  padding: 20px;
  border: 1px solid #d3163b;
  background-color: #fdf2f5;
  box-sizing: border-box;
}
pre {
  border: 1px solid #d3163b;
  padding: 10px;
  background-color: #f8f7f8;
  width: 100%;
  white-space: pre-wrap;
  word-wrap: break-word;
  box-sizing: border-box;
}
h1 {
  font-size: 60px;
  text-align: center;
  margin-top: 50px;
  font-family: OffBit;
}
h2 {
  margin-top: 5px;
  margin-bottom: 5px;
  font-size: 20px;
}
h3 {
  margin-top: 0px;
  font-size: 22px;
  margin-bottom: 45px;
  margin-left: 60px;
  margin-right: 60px;
}
description {
  margin-top: 0px;
  font-size: 25px;
  margin-bottom: 45px;
  font-family: OffBit;
}
.similarity-box {
  margin-bottom: 20px;
}
.text-container {
  display: block;
}
.highlight {
  background-color: yellow;
  color: #d3163b;
}
</style>
</head>
<body>
<h1>Plagiarism Detector •*.✸ </h1>
<h3 style="text-align: center;">Web App in Golang that compares the content of different texts about the same topic, calculating and highlighting their similarity using a Suffix Array that retrieves the Longest Common Substring (LCS).</h3>`

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
            lcs, _ := longestCommonSubstring(text1, text2)
            highlightedText1 := highlightCommonSubstring(text1, lcs)
            highlightedText2 := highlightCommonSubstring(text2, lcs)
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