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
	// Unimos las dos cadenas separadas por un símbolo único
	combined := text1 + "#" + text2 + "$"
	suffixArray := buildSuffixArray(combined)

	n := len(combined)
	longest := ""
	for i := 0; i < n-1; i++ {
		// Verificar si el sufijo pertenece a diferentes textos
		if (suffixArray[i] < len(text1) && suffixArray[i+1] > len(text1)) ||
			(suffixArray[i] > len(text1) && suffixArray[i+1] < len(text1)) {

			// Comparar los sufijos y encontrar la subcadena común más larga
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
	highlighted := strings.Replace(text, commonSubstring, `<span style="background-color: yellow;">`+commonSubstring+`</span>`, -1)
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
<title>Detección de Plagio</title>
<style>
body { font-family: Arial, sans-serif; }
pre { border: 1px solid #000; padding: 10px; background-color: #f4f4f4; }
.highlight { background-color: yellow; }
</style>
</head>
<body>
<h1>Detección de Plagio - Comparación de Textos</h1>`

	for _, pair := range pairs {
		html += fmt.Sprintf("<h2>%s vs %s</h2><pre>%s</pre>", filepath.Base(pair.File1), filepath.Base(pair.File2), pair.Highlighted)
	}

	html += "</body></html>"

	// Guardar el archivo HTML
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

	// Calcular la similitud entre cada par de textos
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			text1 := texts[keys[i]]
			text2 := texts[keys[j]]

			// Encontrar la subcadena común más larga
			lcs, lcsLength := longestCommonSubstring(text1, text2)

			// Calcular la similitud como proporción
			longestLength := max(len(text1), len(text2))
			similarity := float64(lcsLength) / float64(longestLength)

			// Resaltar la subcadena en ambos textos
			highlightedText1 := highlightCommonSubstring(text1, lcs)
			highlightedText2 := highlightCommonSubstring(text2, lcs)

			// Agregar la pareja con su similitud
			pairs = append(pairs, TextPair{
				File1:       keys[i],
				File2:       keys[j],
				Similarity:  similarity,
				Highlighted: fmt.Sprintf("<h3>Similitud: %.2f</h3><pre>%s</pre><pre>%s</pre>", similarity, highlightedText1, highlightedText2),
			})
		}
	}

	// Ordenar los pares por similitud
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Similarity > pairs[j].Similarity
	})

	// Tomar solo los 10 más similares
	if len(pairs) > 10 {
		pairs = pairs[:10]
	}

	// Generar el archivo HTML con los textos resaltados
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
	// Leer archivos desde la carpeta 'dataset'
	texts, err := readFilesFromDir("dataset")
	if err != nil {
		log.Fatal(err)
	}

	// Mostrar los nombres de los archivos leídos
	for fileName, content := range texts {
		fmt.Printf("Archivo: %s\nContenido: %.50s...\n", fileName, content) // Solo muestra los primeros 50 caracteres
	}

	// Calcular la similitud y generar el reporte en HTML
	calculateSimilarityAndHighlight(texts)
}
