class Model:
    '''Модель для поиска'''

    def __init__(self, title: str):
        self.title = str(title)
        self.cursor = 0
        self.len_input = len(str(title))

    def get_current_input(self) -> str:
        return self.title[self.cursor]

    def get_last_input(self) -> str:
        return self.title[self.cursor - 1]

    def increase_cursor(self) -> None:
        self.cursor += 1

    def get_len_input(self) -> int:
        return self.len_input - self.cursor

    def __str__(self):
        return f"model's title = {self.title}"
