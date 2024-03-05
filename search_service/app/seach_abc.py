import pymorphy2


class SearchABC:
    def __init__(self):
        self.unique_words = set()
        self.morph = pymorphy2.MorphAnalyzer()
        self.words_index = {}
        self.bk_trie = None
        self.words_trie = None

    def _words_search(self, string_search: str, cnt_typos: float) -> dict:
        '''Получаем поисковой запрос, отдаем индексы товаров'''
        indexes = {}
        words = string_search.split()
        for word in words:
            word = self.morph.parse(word)[0].normal_form

            if word in self.unique_words:
                searched_models = self.words_trie.predict(word)
                if searched_models:
                    indexes |= self._get_indexes(searched_models, indexes)
            elif not word.isdigit():
                possible_words = self.bk_trie.find(word, int(len(word) * cnt_typos))
                for cnt, possible_word in possible_words:
                    if cnt > 0:
                        searched_models = self.words_trie.predict(possible_word)
                        indexes |= self._get_indexes(searched_models, indexes)

        return indexes

    def _get_indexes(self, searched_models: list, indexes: dict) -> dict:
        for searched_model in searched_models:
            for uuid, vendor_code in self.words_index[searched_model.title]:
                if uuid in indexes:
                    indexes[uuid] += 1
                else:
                    indexes[uuid] = 1
        return indexes
