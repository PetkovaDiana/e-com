import json

from requests import get

url = 'http://go-api:8080/api/v1/products'


def get_products(
        uuids: str,
        price_min: float = 0.1,
        price_max: float = 10000000000.0,
        rating_min: float = 0,
        rating_max: float = 5,
        sort: str = 'default',
        not_empty: str = 'true'
) -> json:
    '''Получение информации о товарах от golang_api'''
    response = get(
        url=url,
        params={
            'prod_id': uuids,
            'price_min': price_min,
            'price_max': price_max,
            'rating_min': rating_min,
            'rating_max': rating_max,
            'sort': sort,
            'not_empty': not_empty
        }
    )
    return response.json()