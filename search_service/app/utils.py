import re


def clean_str(title: str) -> str:
    '''Заменим все спец символы на пробелы'''
    title = str(title)
    title = title.lower()
    answer = ''
    last_char = title[0]
    for char in title:
        if char.isalnum():
            if last_char.isnumeric() != char.isnumeric():
                answer += ' ' + char
            else:
                answer += char
            last_char = char
        elif char == ' ':
            answer += char
        elif char in ['-', '(', '/', '|', '\\', '[', '_', '.', ',']:
            answer += ' '

    answer = re.sub(' +', ' ', answer)
    return answer


def damerau_levenshtein_distance(s1, s2):
    d = {}
    lenstr1, lenstr2 = len(s1), len(s2)
    for i in range(-1, lenstr1 + 1):
        d[(i, -1)] = i + 1
    for j in range(-1, lenstr2 + 1):
        d[(-1, j)] = j + 1

    for i in range(lenstr1):
        for j in range(lenstr2):
            if s1[i] == s2[j]:
                cost = 0
            else:
                cost = 1
            d[(i, j)] = min(
                d[(i - 1, j)] + 1,  # deletion
                d[(i, j - 1)] + 1,  # insertion
                d[(i - 1, j - 1)] + cost,  # substitution
            )
            if i and j and s1[i] == s2[j - 1] and s1[i - 1] == s2[j]:
                d[(i, j)] = min(d[(i, j)], d[i - 2, j - 2] + cost)  # transposition

    return d[lenstr1 - 1, lenstr2 - 1]
