from .model import Model


class Product(Model):
    '''Товар
    Pandas(
    Index='b43e33c5-3e35-11ed-8fd3-002590fb0200',
    title='гирлянда ',
    vendor_code='lgdu121-1-100-10-t-s-44',
    original_title='гирлянда ',
    words=['гирлянда'])
    '''

    def __init__(self, title: str, original_title='', uuid='', vendor_code='') -> None:
        super().__init__(title)
        self.original_title = original_title
        self.uuid = uuid
        self.vendor_code = vendor_code
