from .utils import decimal2text
import decimal


def beautiful(number: float) -> str:
    '''Отступы между тысячами'''
    return '{0:,}'.format(number).replace(',', ' ').replace('.', ',')


def letter_price(price: float) -> str:
    '''Числа прописью'''
    int_units = ((u'рубль', u'рубля', u'рублей'), 'm')
    exp_units = ((u'копейка', u'копейки', u'копеек'), 'f')
    text = decimal2text(
        decimal.Decimal(str(price)),
        int_units=int_units,
        exp_units=exp_units)
    return text