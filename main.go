package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "path/filepath"
    "sort"
    "strings"
)

// Función para leer los archivos de texto desde una carpeta
func readFilesFromDir(dirPath string) (map[string]string, error) {
    files, err := filepath.Glob(dirPath + "/*.txt") // Asume que los archivos son .txt
    if err != nil {
        return nil, err
    }
    texts := make(map[string]string)
    for _, file := range files {
        content, err := ioutil.ReadFile(file)
        if err != nil {
            log.Fatalf("Error leyendo archivo %s: %v", file, err)
        }
        texts[file] = string(content)
    }
    return texts, nil
}

// Función para construir el Suffix Array de una cadena
func buildSuffixArray(text string) []int {
    n := len(text)
    suffixes := make([]string, n)
    suffixArray := make([]int, n)
    for i := 0; i < n; i++ {
        suffixes[i] = text[i:] // Guardamos los sufijos
    }
    // Ordenar los sufijos
    sort.Strings(suffixes)
    // Guardar los índices de los sufijos ordenados
    for i := 0; i < n; i++ {
        suffixArray[i] = n - len(suffixes[i])
    }
    return suffixArray
}

// Función para encontrar la Longest Common Substring (LCS)
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

// Función auxiliar para encontrar el prefijo común más largo entre dos cadenas
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

// Función para resaltar la subcadena común más larga con HTML y CSS
func highlightCommonSubstring(text, commonSubstring string) string {
    highlighted := strings.Replace(text, commonSubstring, `<span style="background-color: yellow; color: #d3163b;">`+commonSubstring+`</span>`, -1)
    return highlighted
}

// Estructura para almacenar la similitud y los pares de textos
type TextPair struct {
    File1       string
    File2       string
    Similarity  float64
    Highlighted string
}

// Función para generar el archivo HTML que muestre los textos con las subcadenas comunes resaltadas
func generateHTML(pairs []TextPair) {
    html := `<html>
<head>
<title>Plagiarism Detector</title>
<style>
body { 
    background-color: #ffa1d7; 
    font-family: Helvetica, sans-serif; 
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
    background-color: #ffffff; 
    box-sizing: border-box; /* Ensure padding fits within width */
}

pre { 
    border: 1px solid #d3163b; 
    padding: 10px; 
    background-color: #f4f4f4; 
    width: 100%; 
    white-space: pre-wrap; 
    word-wrap: break-word;
    box-sizing: border-box;
}

h2 { 
    margin-bottom: 5px; 
    font-size: 20px;
}

.similarity-box {
    margin-bottom: 20px;
}

.text-container {
    display: block; /* Stack the contents vertically */
}

.highlight { 
    background-color: yellow; 
    color: #d3163b; 
}
</style>
</head>
<body>
<h1 style="text-align: center;">Plagiarism Detector</h1>
<h3 style="text-align: center;">Web App in Golang that compares the content of different texts about the same topic, calculating and highlighting their similarity using a Suffix Array that retrieves the Longest Common Substring (LCS).</h3>`

    for _, pair := range pairs {
        // Split the highlighted string into two parts for each text
        splitIndex := strings.Index(pair.Highlighted, "</pre><pre>") + 6
        firstText := pair.Highlighted[:splitIndex]   // Extract the first file's highlighted text
        secondText := pair.Highlighted[splitIndex:]  // Extract the second file's highlighted text

        html += fmt.Sprintf(`
        <div class="container">
            <h2>%s vs %s</h2>
            <h3>Similitud: %.2f</h3> <!-- Similarity is printed here under the title only -->
            <div class="text-container">
                <div class="text-content">
                    <pre>%s</pre> <!-- Only text content is printed here -->
                </div>
                <div class="text-content">
                    <pre>%s</pre>
                </div>
            </div>
        </div>`, 
            filepath.Base(pair.File1), filepath.Base(pair.File2), pair.Similarity, firstText, secondText)
    }

    html += "</body></html>"
    
    err := ioutil.WriteFile("plagiarism_report.html", []byte(html), 0644)
    if err != nil {
        log.Fatalf("Error al escribir el archivo HTML: %v", err)
    }
    fmt.Println("Archivo HTML generado: plagiarism_report.html")
}

// Función para calcular la similitud y generar los pares de texto
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
            lcs, lcsLength := longestCommonSubstring(text1, text2)
            longestLength := max(len(text1), len(text2))
            similarity := float64(lcsLength) / float64(longestLength)
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
    sort.Slice(pairs, func(i, j int) bool {
        return pairs[i].Similarity > pairs[j].Similarity
    })
    if len(pairs) > 10 {
        pairs = pairs[:10]
    }
    generateHTML(pairs)
}

// Función auxiliar para obtener el máximo de dos enteros
func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func main() {
    texts, err := readFilesFromDir("dataset")
    if err != nil {
        log.Fatal(err)
    }
    for fileName, content := range texts {
        fmt.Printf("Archivo: %s\nContenido: %.50s...\n", fileName, content) // Solo muestra los primeros 50 caracteres
    }
    calculateSimilarityAndHighlight(texts)
}
