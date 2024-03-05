import pandas as pd

from .utils import clean_str, damerau_levenshtein_distance
from .tries.bk import BKTree
from .tries.jobs_trie import JobsTrie, WordsJobTrie
from .seach_abc import SearchABC


class JobsSearch(SearchABC):
    '''Поиск по вакансиям'''

    def __init__(self, df: pd.DataFrame):
        super().__init__()
        self.jobs = df
        self.jobs_trie = JobsTrie()
        self.words_trie = WordsJobTrie()

    def __preparation(self):
        self.jobs['title'] = self.jobs['title'].str.lower()
        self.jobs = self.jobs.drop_duplicates(subset=['title'])
        self.jobs['original_title'] = self.jobs['title']
        self.jobs['title'] = self.jobs['title'].apply(clean_str)
        self.jobs['words'] = self.jobs['title'].str.split()

        unique_words_all = set()
        for product_words in self.jobs['words']:
            unique_words_all.update(product_words)

        for unique_word in unique_words_all:
            if len(unique_word) > 1:
                self.unique_words.add(self.morph.parse(unique_word)[0].normal_form)

        for row in self.jobs.itertuples():
            for word in row[-1]:
                if len(word) > 1:
                    word = self.morph.parse(word)[0].normal_form
                    if word in self.words_index:
                        self.words_index[word].add(row[0])
                    else:
                        self.words_index[word] = {row[0]}

        self.bk_trie = BKTree(damerau_levenshtein_distance, self.unique_words)

    def fit(self):
        self.__preparation()
        self.jobs_trie.fit_on_jobs(self.jobs)
        self.words_trie.fit_on_words(self.words_index)

    def _get_indexes(self, searched_models: list, indexes: dict) -> dict:
        for searched_model in searched_models:
            for id in self.words_index[searched_model.title]:
                if id in indexes:
                    indexes[id] += 1
                else:
                    indexes[id] = 1
        return indexes
