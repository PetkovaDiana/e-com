from .model import Model


class Node:
    '''Узел префиксного дерева'''

    def __init__(self):
        self.keys = {}
        self.end = False
        self.model = None

    def _set_end(self):
        self.end = True

    def set_model(self, model: Model):
        self._set_end()
        self.model = model

    def is_end(self):
        return self.end

    def __str__(self):
        return f'keys = {self.keys}'