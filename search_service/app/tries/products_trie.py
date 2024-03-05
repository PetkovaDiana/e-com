import pandas as pd

from .trie import Trie
from ..models.product import Product


class ProductsTrie(Trie):
    '''Префиксное дерево полных названий товаров'''

    def fit_on_products(self, products: pd.DataFrame) -> None:
        for row in products.itertuples():
            product = Product(
                title=row[1],
                original_title=row[3],
                uuid=row[0],
                vendor_code=row[2]
            )
            self.add(product, self.root)


class WordsTrie(Trie):
    '''Префиксное дерево слов из названий товаров'''

    def fit_on_words(self, word_index: dict) -> None:
        '''Префиксное дерево на словах в составе названий'''
        for word in word_index:
            self.add(Product(word), self.root)


class VendorCodesTrie(Trie):
    '''Префиксное дерево на артикулах'''

    def fit_on_vendor_codes(self, products: pd.DataFrame) -> None:
        for row in products.itertuples():
            product = Product(
                title=row[2],
                original_title=row[3],
                uuid=row[0],
                vendor_code=row[1]
            )
            self.add(product, self.root)