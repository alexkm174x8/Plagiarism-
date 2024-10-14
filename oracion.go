package main

import (
	"fmt"
	"sort"
	"strings"
)

// Estructura para guardar los sufijos (listas de palabras) y sus índices
type Suffix struct {
	suffix []string
	index  int
}

// Función para construir el Suffix Array basado en palabras
func buildSuffixArrayWords(words []string) []int {
	var suffixes []Suffix
	for i := range words {
		suffixes = append(suffixes, Suffix{suffix: words[i:], index: i})
	}
	sort.Slice(suffixes, func(i, j int) bool {
		// Comparar los sufijos palabra por palabra
		return strings.Join(suffixes[i].suffix, " ") < strings.Join(suffixes[j].suffix, " ")
	})

	var suffixArray []int
	for _, suf := range suffixes {
		suffixArray = append(suffixArray, suf.index)
	}
	return suffixArray
}

// Función para encontrar la subcadena común más larga en palabras
func longestCommonSubstringWords(documents []string) []string {
	var commonSubstrings []string
	for i := 0; i < len(documents); i++ {
		for j := i + 1; j < len(documents); j++ {
			doc1Words := strings.Split(documents[i], " ")
			doc2Words := strings.Split(documents[j], " ")
			lcsLength := 0
			longestCommon := ""

			// Construir suffix array por palabras para ambos documentos
			suffixArray1 := buildSuffixArrayWords(doc1Words)
			suffixArray2 := buildSuffixArrayWords(doc2Words)

			// Comparar sufijos para encontrar la coincidencia más larga
			for _, idx1 := range suffixArray1 {
				for _, idx2 := range suffixArray2 {
					k := 0
					for idx1+k < len(doc1Words) && idx2+k < len(doc2Words) && doc1Words[idx1+k] == doc2Words[idx2+k] {
						k++
					}
					if k > lcsLength {
						lcsLength = k
						longestCommon = strings.Join(doc1Words[idx1:idx1+k], " ")
					}
				}
			}

			if longestCommon == "" {
				longestCommon = "No common substring"
			}
			commonSubstrings = append(commonSubstrings, fmt.Sprintf("Doc %d vs Doc %d - Subcadena común más larga: '%s' (longitud: %d palabras)", i+1, j+1, longestCommon, lcsLength))
		}
	}
	return commonSubstrings
}

func main() {
	documents := []string{
		"el gato saltó sobre la mesa",
		"el perro corrió por el parque",
		"el gato y el perro son amigos",
		"la mesa estaba cerca de la puerta",
		"el sol brilla en el cielo azul",
		"las nubes cubren el sol a veces",
		"el parque está lleno de gente hoy",
		"los niños juegan en el parque soleado",
		"el perro corre detrás del gato",
		"los amigos se sientan en la mesa",
	}

	results := longestCommonSubstringWords(documents)
	for _, result := range results {
		fmt.Println(result)
	}
}
