from collections import deque
from operator import itemgetter

_getitem0 = itemgetter(0)


class BKTree(object):
    """
    BK-древовидная структура данных, позволяющая быстро запрашивать совпадения, которые являются
    "закрыть" задана функция для вычисления показателя расстояния (например, Хэмминг
    расстояние или расстояние Левенштейна).

    Каждый узел в дереве (включая корневой узел) представляет собой кортеж из двух
    элементов (item, children_dict), где children_dict - это dict, ключами которого являются
    неотрицательные расстояния дочернего элемента до текущего элемента и значения которого
    являются узлами.
    """

    def __init__(self, distance_func, items):
        """Инициализируйте экземпляр BKTree с заданной функцией расстояния
        (который принимает два элемента в качестве параметров и возвращает неотрицательное значение
        целое число расстояний). "элементы" - это необязательный список элементов, добавляемых
        при инициализации. """
        self.distance_func = distance_func
        self.tree = None

        _add = self.add
        for item in items:
            _add(item)

    def add(self, item):
        """Добавьте данный элемент в это дерево. """
        node = self.tree
        if node is None:
            self.tree = (item, {})
            return

        # Небольшая оптимизация скорости - избегайте поиска внутри цикла
        _distance_func = self.distance_func

        while True:
            parent, children = node
            distance = _distance_func(item, parent)
            node = children.get(distance)
            if node is None:
                children[distance] = (item, {})
                break

    def find(self, item, n):
        """Найдите элементы в этом дереве, расстояние между которыми меньше или равно n
        из заданного элемента и возвращает список кортежей (расстояние, элемент), упорядоченных по
        расстоянию. """
        if self.tree is None:
            return []

        candidates = deque([self.tree])
        found = []

        # Небольшая оптимизация скорости - избегайте поиска внутри цикла
        _candidates_popleft = candidates.popleft
        _candidates_extend = candidates.extend
        _found_append = found.append
        _distance_func = self.distance_func

        while candidates:
            candidate, children = _candidates_popleft()
            distance = _distance_func(candidate, item)
            if distance <= n:
                _found_append((distance, candidate))

            if children:
                lower = distance - n
                upper = distance + n
                _candidates_extend(c for d, c in children.items() if lower <= d <= upper)

        found.sort(key=_getitem0)
        return found

    def __iter__(self):
        """Возвращает итератор по всем элементам в этом дереве; элементы выводятся в
        произвольном порядке. """
        if self.tree is None:
            return

        candidates = deque([self.tree])

        # Небольшая оптимизация скорости - избегайте поиска внутри цикла
        _candidates_popleft = candidates.popleft
        _candidates_extend = candidates.extend

        while candidates:
            candidate, children = _candidates_popleft()
            yield candidate
            _candidates_extend(children.values())

    def __repr__(self):
        """Верните строковое представление этого BK-дерева с небольшим количеством информации."""
        return '<{} using {} with {} top-level nodes>'.format(
            self.__class__.__name__,
            self.distance_func.__name__,
            len(self.tree[1]) if self.tree is not None else 'no',
        )
