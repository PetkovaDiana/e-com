from .model import Model


class Job(Model):
    '''Вакансия
    Pandas(
    0   Index=3,
    1   first_phone='+7 (908) 355-28-88',
    2   second_phone=None,
    3   title='водитель',
    4   email='sales@besmufa.ru',
    5   original_title='водитель',
    6   words=['водитель']
    )
'''

    def __init__(self, title: str, original_title='', id=0, first_phone='', second_phone='', email=''):
        super().__init__(title)
        self.id = id
        self.first_phone = first_phone
        self.second_phone = second_phone
        self.email = email
        self.original_title = original_title
