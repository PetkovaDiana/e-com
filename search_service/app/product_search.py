import pandas as pd

from .utils import clean_str, damerau_levenshtein_distance
from .tries.bk import BKTree
from .tries.products_trie import ProductsTrie, WordsTrie, VendorCodesTrie
from .seach_abc import SearchABC


class ProductSearch(SearchABC):
    '''Алгоритм поиска'''

    def __init__(self, df: pd.DataFrame):
        super().__init__()
        self.products = df
        self.products_trie = ProductsTrie()
        self.words_trie = WordsTrie()
        self.vendor_code_trie = VendorCodesTrie()

    def __preparation(self):
        # Приведем к нижнему регистру
        self.products['title'] = self.products['title'].str.lower()
        self.products['vendor_code'] = self.products['vendor_code'].str.lower()

        # Удалим дубликаты (название и артикул одновременно совпадают)
        self.products = self.products.drop_duplicates(subset=['title', 'vendor_code'])
        self.products['original_title'] = self.products['title']

        # Удалим товары без названия
        self.products = self.products.loc[self.products['title'] != '']

        # Заменим все спец символы на пробелы в названии товаров
        #self.products['title'] = self.products['title'].apply(clean_str) #TODO заказчик попросил учитывать в поиске
        # спец символы

        # Удалим "ект" из названий товаров
        self.products['title'] = self.products['title'].str.replace(' ект', '', regex=True)
        self.products['title'] = self.products['title'].str.replace('иэк', 'iek', regex=True)

        # Разобьем названия товаров на слова
        self.products['words'] = self.products['title'].str.split()

        # Создание обратного индексированияю. Составим множество уникальный слов
        unique_words_all = set()
        for product_words in self.products['words']:
            if type(product_words) != float:
                unique_words_all.update(product_words)

        #Приведем слова в начальную форму
        for unique_word in unique_words_all:
            # if len(unique_word) > 1:
            self.unique_words.add(self.morph.parse(unique_word)[0].normal_form)

        #Исключим часто употребимые слова которые не несут смысловой нагрузки
        common_words = ['для', 'шт', 'на', 'по', 'мм']
        for common_word in common_words:
            if common_word in self.unique_words:
                self.unique_words.remove(common_word)

        #Построим обратное индексирование
        for row in self.products.itertuples():
            if type(row[-1]) == float:
                continue
            for word in row[-1]:
                if len(word) > 1:
                    word = self.morph.parse(word)[0].normal_form
                    if word in self.words_index:
                        self.words_index[word].add((row[0], row[2]))
                    else:
                        self.words_index[word] = {(row[0], row[2])}

        self.bk_trie = BKTree(damerau_levenshtein_distance, self.unique_words)

    def fit(self):
        print(f"SearchService: ProductSearch.fit(): start")
        self.__preparation()
        print('SearchService: ProductSearch.fit(): Начали заполнение префиксного дерева поиска по названию')
        self.products_trie.fit_on_products(self.products)
        print('SearchService: ProductSearch.fit(): Начали заполнение префиксного дерева поиска по названию')
        print('SearchService: ProductSearch.fit(): Начали заполнение префиксного дерева поиска по словам в названии')
        self.words_trie.fit_on_words(self.words_index)
        print('SearchService: ProductSearch.fit(): Закончили заполнение префиксного дерева поиска по словам в названии')
        print('SearchService: ProductSearch.fit(): Начали заполнение префиксного дерева поиска по артикулу')
        self.vendor_code_trie.fit_on_vendor_codes(self.products)
        print('SearchService: ProductSearch.fit(): Закончили заполнение префиксного дерева поиска по артикулу')