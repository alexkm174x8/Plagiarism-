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
// Lee todos los archivos .txt en un directorio y devuelve un mapa con el nombre del archivo y su contenido
func readFilesFromDir(dirPath string) (map[string]string, error) {
	// Obtenemos la lista de archivos .txt en el directorio
	files, err := filepath.Glob(dirPath + "/*.txt")
	if err != nil {
		return nil, err
	}
	texts := make(map[string]string)

	// Leemos cada archivo y almacenamos su contenido en el mapa
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
// Un Suffix Array es un arreglo que contiene los índices de los sufijos de la cadena en orden lexicográfico
func buildSuffixArray(text string) []int {
	n := len(text)
	suffixes := make([]string, n) // Arreglo de sufijos
	suffixArray := make([]int, n) // Arreglo para almacenar el índice del sufijo

	// Construimos los sufijos
	for i := 0; i < n; i++ {
		suffixes[i] = text[i:] // Cada sufijo comienza desde la posición i
	}

	// Ordenamos los sufijos lexicográficamente
	sort.Strings(suffixes)

	// Guardamos los índices de los sufijos en el arreglo de Suffix Array
	for i := 0; i < n; i++ {
		suffixArray[i] = n - len(suffixes[i])
	}

	return suffixArray
}

// Función para encontrar la Longest Common Substring (LCS)
// Compara dos cadenas y encuentra la subcadena común más larga usando el Suffix Array
func longestCommonSubstring(text1, text2 string) (string, int) {
	// Combinamos ambas cadenas con un delimitador (# y $ para evitar colisiones)
	combined := text1 + "#" + text2 + "$"

	// Construimos el Suffix Array para la cadena combinada
	suffixArray := buildSuffixArray(combined)

	n := len(combined)
	longest := ""

	// Iteramos sobre el Suffix Array buscando sufijos de ambas cadenas
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
// Compara las dos cadenas carácter por carácter y devuelve el prefijo común
func commonPrefix(s1, s2 string) string {
	n := len(s1)
	if len(s2) < n {
		n = len(s2)
	}

	// Iteramos sobre los caracteres de ambas cadenas
	for i := 0; i < n; i++ {
		if s1[i] != s2[i] {
			return s1[:i] // Devolvemos la subcadena común hasta el punto de diferencia
		}
	}
	return s1[:n]
}

// Función para resaltar la subcadena común más larga con HTML y CSS
// Reemplaza la subcadena común con una versión resaltada usando HTML
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
// Genera un archivo HTML con los textos comparados y resaltados en amarillo donde coinciden
func generateHTML(pairs []TextPair) {
	// Plantilla del HTML con estilos CSS
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

	// Iteramos sobre cada par de textos y agregamos la comparación al HTML
	for _, pair := range pairs {
		splitIndex := strings.Index(pair.Highlighted, "</pre><pre>") + 6
		firstText := pair.Highlighted[:splitIndex]
		secondText := pair.Highlighted[splitIndex:]

		html += fmt.Sprintf(`
        <div class="container">
            <description>%s vs %s</description>
            <h2>Similarity: %.2f</h2>
            <div class="text-container">
                <div class="text-content">
                    <pre>%s</pre>
                </div>
                <div class="text-content">
                    <pre>%s</pre>
                </div>
            </div>
        </div>`,
			filepath.Base(pair.File1), filepath.Base(pair.File2), pair.Similarity, firstText, secondText)
	}

	html += "</body></html>"

	// Guardamos el HTML en un archivo
	err := ioutil.WriteFile("plagiarism_report.html", []byte(html), 0644)
	if err != nil {
		log.Fatalf("Error al escribir el archivo HTML: %v", err)
	}
	fmt.Println("Archivo HTML generado: plagiarism_report.html")
}

// Función para calcular la similitud y generar los pares de texto
// Compara todos los textos entre sí y genera una matriz de similitud con los 10 pares más similares
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
	// Ordenamos los pares por similitud y seleccionamos los 10 más similares
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
	// Leemos los textos del directorio 'dataset'
	texts, err := readFilesFromDir("dataset")
	if err != nil {
		log.Fatal(err)
	}

	// Mostramos una vista previa de los archivos leídos
	for fileName, content := range texts {
		fmt.Printf("Archivo: %s\nContenido: %.50s...\n", fileName, content) // Solo muestra los primeros 50 caracteres
	}

	// Calculamos la similitud entre los textos y generamos el archivo HTML
	calculateSimilarityAndHighlight(texts)
}
