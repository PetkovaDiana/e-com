from fastapi import APIRouter

from .crud import RecommendationApi
from .db_connect import PandasDataBase


router = APIRouter()
db: PandasDataBase = PandasDataBase()

product_recommendation_api = RecommendationApi(db.get_products())
product_recommendation_api.fit()


@router.put('/recommendation/')
async def put_product(products_uuids: str):
    '''
    Обновление весов, в случае продажи
    вернет список неизвестных вершин, в идеале это []
    '''
    return product_recommendation_api.update(products_uuids.split(','))


@router.get('/recommendation/')
async def get_product_recommendation(product_uuid: str, limit: int = 10):
    '''Получение списка рекомендованных товаров'''
    return product_recommendation_api.predict(product_uuid, limit)