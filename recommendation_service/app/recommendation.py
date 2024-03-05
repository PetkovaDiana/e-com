import pandas as pd
import numpy as np
import json


class Recommendation:
    '''
    С этим товаром часто покупают, работа с взвешенным графом
    verb_* - id товара
    структура - граф (список смежности)
    '''

    def __init__(self, products: pd.DataFrame):
        self.products_uuid_array: np.array = products.index.to_numpy()
        self.data: dict = dict()

    def fit(self) -> None:
        '''Заполнение нашего графа'''
        with open("app/data.json", "r") as read_file:
            data = read_file
        if not self.data:
            for verb_uuid in self.products_uuid_array:
                self.data[verb_uuid] = {}
        else:
            self.data = json.load(data)

    def update(self, verbs_uuid: list) -> list:
        '''
        Обновление весов, в случае продажи
        вернет список неизвестных вершин, в идеале это []
        '''
        unknown_uuids = []
        for verb_uuid in verbs_uuid:
            if verb_uuid in self.data:
                for second_uuid in verbs_uuid:
                    if verb_uuid != second_uuid:
                        if second_uuid in self.data[verb_uuid]:
                            self.data[verb_uuid][second_uuid] += 1
                        else:
                            self.data[verb_uuid][second_uuid] = 1
            else:
                unknown_uuids.append(verb_uuid)

        with open('app/data.json', 'w') as outfile:
            json.dump(self.data, outfile, ensure_ascii=False)
        return unknown_uuids

    def get(self, verb_uuid: str, limit: int = 10) -> list:
        '''Получение списка рекомендованных товаров'''
        if verb_uuid in self.data:
            verbs: dict = self.data[verb_uuid]
            sorted_tuple = sorted(verbs.items(), key=lambda x: x[1], reverse=True)
            if sorted_tuple:
                sorted_list = list(zip(*sorted_tuple))[0][:limit]
            else:
                sorted_list = []
            return sorted_list
        else:
            return []