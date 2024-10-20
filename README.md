# Sistema de Detección de Plagio

El plagio ha aumentado con el acceso masivo a información en línea, y la identificación de plagio en grandes volúmenes de texto se ha vuelto una tarea compleja. Este proyecto aborda el problema implementando un sistema de detección de similitudes textuales utilizando **Go** y algoritmos de procesamiento de cadenas, como el **Suffix Array** y el **Longest Common Substring (LCS)**.

Este sistema identifica las subcadenas más largas comunes entre pares de textos y destaca visualmente las coincidencias en un reporte HTML. El proyecto también explora las complejidades computacionales, inspirándose en trabajos relevantes como "All Optimal k-Bounded Alignments Using the FM-Index" y "Developing a corpus of plagiarised short answers".

### Backend

El sistema se basa en el **Suffix Array** para procesar y comparar subcadenas. Las principales fases incluyen:

- **Lectura de archivos**: Se leen archivos de texto desde un directorio específico.
- **Construcción del Suffix Array**: Se utiliza el algoritmo **Merge Sort** para ordenar los sufijos lexicográficamente.
- **Cálculo de la Longest Common Substring (LCS)**: Se comparan los sufijos de dos textos para encontrar las subcadenas comunes más largas.
- **Eliminación de subcadenas comunes**: Se eliminan iterativamente las subcadenas comunes hasta un umbral.
- **Cálculo de similitud mejorado**: Se usa el promedio de las longitudes originales de los textos para mantener la similitud en un rango válido.
- **Ordenación de pares por similitud**: Los resultados se ordenan según la similitud, facilitando el análisis.

### Frontend

Se genera un reporte HTML que resalta visualmente las similitudes entre textos, lo que facilita la identificación de plagio.

## Análisis de Complejidad

- **Construcción del Suffix Array**: O(n log n)
- **Eliminación de subcadenas comunes**: O(n² log n)

Aunque eficiente en textos moderados, se exploran optimizaciones como el uso de **FM-Index** o **hashing** para reducir la complejidad en grandes volúmenes de datos.

## Resultados

El sistema fue capaz de detectar similitudes entre textos de forma precisa, destacando las áreas de coincidencia con una visualización clara. El tiempo de procesamiento fue optimizado, logrando resultados comparables a soluciones que usan el **FM-Index**.

## Contribuciones de los Integrantes

### Alejandro Kong Montoya

Alejandro trabajó en la **implementación base del algoritmo** y en la **construcción del Suffix Array**. Además, fue responsable de la **generación del reporte HTML** para la visualización de los resultados. También se encargó de **mejorar el código** para optimizar el rendimiento del sistema.

### Estefanía Antonio Villaseca

Estefanía desarrolló la **base del código para la implementación del Longest Common Substring (LCS)**. También refinó el **cálculo de la similitud** entre los textos, asegurando que el sistema identificara las coincidencias con precisión y eficiencia.

### Miranda Eugenia Colorado Arróniz

Miranda trabajó en la **presentación de los resultados** de la comparación de textos en el frontend del reporte HTML. Además, diseñó las **funciones en Golang para el preprocesamiento** y la **eliminación de subcadenas**, optimizando los resultados generales del programa y mejorando su precisión.

## Reflexiones

### Estefanía Antonio Villaseca

Este proyecto me ayudó a desarrollar habilidades en el uso de estructuras de datos avanzadas y a aplicar buenas prácticas de programación para asegurar un código escalable y mantenible.

### Miranda Eugenia Colorado Arróniz

Me enfoqué en mejorar la experiencia del usuario y optimizar el cálculo de similitud, incluyendo preprocesamiento de texto y visualización de resultados.

### Alejandro Kong Montoya

Logramos implementar una solución eficiente que destaca visualmente las coincidencias en un informe HTML. En el futuro, se podrían integrar algoritmos adicionales como el **FM-Index** para mejorar el rendimiento.

## Conclusión

Este proyecto ofrece una solución eficiente para la detección de plagio, combinando algoritmos de procesamiento de cadenas y una presentación visual clara. La herramienta es valiosa para ámbitos académicos y profesionales donde el plagio es una preocupación creciente.
