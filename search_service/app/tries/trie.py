from ..models.model import Model
from ..models.node import Node


class Trie:
    '''Префиксное дерево'''

    def __init__(self):
        self.root = Node()

    def _get_model_node(self, title: str) -> Node:
        node = self.root
        while len(title) > 1:
            if title[0] in node.keys:
                node = node.keys[title[0]]
                title = title[1:]
            else:
                return
        return node.keys.get(title[0])

    def add(self, model: Model, node: Node):
        if model.get_len_input() == 0:
            node.set_model(model)
            return
        elif model.get_current_input() not in node.keys:
            node.keys[model.get_current_input()] = Node()
            model.increase_cursor()
            return self.add(model, node.keys[model.get_last_input()])
        else:
            model.increase_cursor()
            return self.add(model, node.keys[model.get_last_input()])

    def is_product(self, title: str):
        '''Является ли данная строка названием товара'''
        node = self.root

        while len(title) > 1:
            if title[0] not in node.keys:
                return False
            else:
                node = node.keys[title[0]]
                title = title[1:]
        if title in node.keys and node.keys[title].is_end():
            return node.keys[title]
        else:
            return False

    def _get_models(self, node: Node):
        def __search(node: Node, string: str) -> None:
            for letter in node.keys:
                yield from __search(node.keys[letter], string + letter)

            if node.is_end():
                yield node.model

        yield from __search(node, '')

    def predict(self, title: str):
        fast_node = self.is_product(title)
        if fast_node:
            yield fast_node.model
        else:
            searched_node = self._get_model_node(title)
            if searched_node:
                yield from self._get_models(searched_node)