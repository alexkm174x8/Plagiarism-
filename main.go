package main


def build_suffix_array(s):
    suffixes = [(s[i:], i) for i in range(len(s))]
    suffixes.sort()  # Ordena los sufijos lexicográficamente
    suffix_array = [suf[1] for suf in suffixes]  # Almacena los índices
    return suffix_array

def build_bwt(s, suffix_array):
    bwt = ''.join([s[i-1] if i > 0 else '$' for i in suffix_array])
    return bwt


def build_c_and_occ(bwt):
    alphabet = sorted(set(bwt))
    C = {}
    Occ = {char: [0] * (len(bwt) + 1) for char in alphabet}
    
    for i in range(1, len(bwt) + 1):
        char = bwt[i-1]
        for a in alphabet:
            Occ[a][i] = Occ[a][i-1] + (1 if char == a else 0)
    
    total = 0
    for a in alphabet:
        C[a] = total
        total += Occ[a][-1]
    
    return C, Occ

def backward_search(pattern, bwt, C, Occ):
    l = 0
    r = len(bwt)
    for char in reversed(pattern):
        if char not in C:
            return -1, -1  # No hay coincidencias
        l = C[char] + Occ[char][l]
        r = C[char] + Occ[char][r]
        if l >= r:
            return -1, -1  # No hay coincidencias
    return l, r

def longest_common_substring(documents):
    concatenated_text = '|'.join(documents) + '$'
    
    suffix_array = build_suffix_array(concatenated_text)
    bwt = build_bwt(concatenated_text, suffix_array)
    
    C, Occ = build_c_and_occ(bwt)
    
    common_substrings = []
    for i in range(len(documents)):
        for j in range(i + 1, len(documents)):
            lcs_length = 0
            doc1 = documents[i]
            doc2 = documents[j]
            for k in range(len(doc1)):
                pattern = doc1[k:]  # Busca subcadenas de doc1
                l, r = backward_search(pattern, bwt, C, Occ)
                if l != -1 and r != -1:
                    # Si hay coincidencias, calculamos la longitud de la subcadena común
                    match_length = r - l
                    if match_length > lcs_length:
                        lcs_length = match_length
            common_substrings.append((f"Doc {i+1} vs Doc {j+1}", lcs_length))
    
    return common_substrings

documents = [
    "el gato saltó sobre la mesa",
    "el perro corrió por el parque",
    "el gato y el perro son amigos",
    "la mesa estaba cerca de la puerta",
    "el sol brilla en el cielo azul",
    "las nubes cubren el sol a veces",
    "el parque está lleno de gente hoy",
    "los niños juegan en el parque soleado",
    "el perro corre detrás del gato",
    "los amigos se sientan en la mesa"
]

results = longest_common_substring(documents)

for result in results:
    print(f"{result[0]} - Longitud de subcadena común más larga: {result[1]}")