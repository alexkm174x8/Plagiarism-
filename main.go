package main

// Importación de paquetes necesarios para la funcionalidad del programa.
// Cada paquete cumple una función esencial en la lectura, manipulación y procesamiento de archivos y cadenas.
import (
	"fmt"           // Proporciona funciones para formatear y mostrar datos en la consola.
	"io/ioutil"     // Permite realizar operaciones de entrada y salida, como leer y escribir archivos.
	"log"           // Facilita la creación de registros de eventos, especialmente para manejar errores.
	"path/filepath" // Proporciona funciones para manipular rutas de archivos de forma independiente del sistema operativo.
	"strings"       // Incluye funciones para la manipulación de cadenas de texto, como la búsqueda, reemplazo y segmentación.
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
        processedContent := preprocessText(string(content))
		texts[file] = processedContent
	}
	return texts, nil
}

// Función para preprocesar el texto: elimina saltos de línea, tabulaciones y retornos de carro, pero mantiene los espacios
func preprocessText(text string) string {
    // Elimina saltos de línea, retornos de carro y tabulaciones
    cleaned := strings.ReplaceAll(text, "\n", "")   // Elimina saltos de línea
    cleaned = strings.ReplaceAll(cleaned, "\r", "") // Elimina retornos de carro
    cleaned = strings.ReplaceAll(cleaned, "\t", "") // Elimina tabulaciones
    cleaned = strings.TrimSpace(cleaned)            // Elimina espacios al inicio y al final (pero conserva los espacios dentro del texto)
    
    return cleaned
}

// Implementación de Merge Sort para ordenar sufijos lexicográficamente
func mergeSortSuffixes(suffixes []string) []string {
	if len(suffixes) <= 1 {
		return suffixes
	}

	mid := len(suffixes) / 2
	left := mergeSortSuffixes(suffixes[:mid])
	right := mergeSortSuffixes(suffixes[mid:])

	return merge(left, right)
}



// Función auxiliar para mezclar dos sublistas en orden lexicográfico
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

	// Añadir los elementos restantes
	result = append(result, left[i:]...)
	result = append(result, right[j:]...)

	return result
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

	// Ordenamos los sufijos lexicográficamente usando Merge Sort
	sortedSuffixes := mergeSortSuffixes(suffixes)

	// Guardamos los índices de los sufijos en el arreglo de Suffix Array
	for i := 0; i < n; i++ {
		suffixArray[i] = n - len(sortedSuffixes[i])
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

// Función para eliminar la subcadena común más larga de ambos textos
// Continúa buscando y eliminando subcadenas comunes más largas de 5 caracteres hasta que ya no se encuentren
func removeLongestSubstrings(text1, text2 string) ([]string, string, string) {
    removedSubstrings := []string{}
    
    // Bucle para encontrar y eliminar subcadenas comunes de longitud mayor a 5
    for {
        lcs, lcsLength := longestCommonSubstring(text1, text2)
        
        // Detener el bucle si la subcadena más larga tiene 5 caracteres o menos
        if lcsLength <= 5 {
            break
        }
        
        // Almacenar la subcadena eliminada
        removedSubstrings = append(removedSubstrings, lcs)
        
        // Eliminar todas las ocurrencias de la subcadena de ambos textos
        text1 = strings.ReplaceAll(text1, lcs, "")
        text2 = strings.ReplaceAll(text2, lcs, "")
    }
    
    return removedSubstrings, text1, text2
}


// Función para resaltar la subcadena común más larga con HTML y CSS
// Reemplaza la subcadena común con una versión resaltada usando HTML
func highlightCommonSubstrings(text string, commonSubstrings []string) string {
    highlighted := text
    for _, substr := range commonSubstrings {
        highlighted = strings.ReplaceAll(highlighted, substr, `<span style="background-color: yellow; color: #d3163b;">`+substr+`</span>`)
    }
    return highlighted
}

// Estructura para almacenar la similitud y los pares de textos
type TextPair struct {
	File1       string
	File2       string
	Similarity  float64
	Highlighted string
}

// Implementación de Merge Sort para ordenar pares de textos según su similitud
func mergeSortPairs(pairs []TextPair) []TextPair {
	if len(pairs) <= 1 {
		return pairs
	}

	mid := len(pairs) / 2
	left := mergeSortPairs(pairs[:mid])
	right := mergeSortPairs(pairs[mid:])

	return mergePairs(left, right)
}

// Función auxiliar para mezclar dos sublistas de TextPair en orden de similitud
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

	// Añadir los elementos restantes
	result = append(result, left[i:]...)
	result = append(result, right[j:]...)

	return result
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

	counter := 1

	// Iteramos sobre cada par de textos y agregamos la comparación al HTML
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

	// Guardamos el HTML en un archivo
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

            // Remove longest substrings iteratively and collect them
            removedSubstrings, finalText1, finalText2 := removeLongestSubstrings(text1, text2)
            
            // Calculate similarity based on removed substrings
            totalCommonLength := 0
            for _, common := range removedSubstrings {
                totalCommonLength += len(common)
            }
            
            longestLength := max(len(finalText1), len(finalText2))
            similarity := float64(totalCommonLength) / float64(longestLength)
            
            highlightedText1 := highlightCommonSubstrings(text1, removedSubstrings)
            highlightedText2 := highlightCommonSubstrings(text2, removedSubstrings)

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