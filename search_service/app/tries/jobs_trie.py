import pandas as pd

from ..models.job import Job
from .trie import Trie


class JobsTrie(Trie):
    '''Префиксное дерево для вакансий'''

    def fit_on_jobs(self, jobs: pd.DataFrame) -> None:
        for row in jobs.itertuples():
            job = Job(
                id=row[0],
                title=row[3],
                original_title=row[5],
                first_phone=row[1],
                second_phone=row[2],
                email=row[4]
            )
            self.add(job, self.root)


class WordsJobTrie(Trie):
    '''Префиксное дерево на словах в составе названий'''

    def fit_on_words(self, word_index: dict) -> None:
        for word in word_index:
            self.add(Job(word), self.root)
