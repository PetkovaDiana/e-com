from fastapi import APIRouter
from fastapi_pagination import paginate, Page

from .crud import SearchJobsApi, SearchProductsApi
from .db.db_connect import PandasDataBase
from .db.db_models import SearchWindow, SearchJob, Products

router = APIRouter()
db: PandasDataBase = PandasDataBase()

product_search_api = SearchProductsApi(db.get_products())
product_search_api.fit()

job_search_api = SearchJobsApi(db.get_job())
job_search_api.fit()


@router.get("/search/", response_model=SearchWindow)
async def get_search(search_string: str, limit: int = 10):
    '''
    Поиск по названию товара
    '''
    answer: list = product_search_api.predict(search_string, limit, .2)
    return {'data': answer}


@router.get("/search_page/", response_model=Page[Products])
async def get_search(
        search_string: str,
        limit: int = 10,
        price_min: float = 0.1,
        price_max: float = 10000000000.0,
        rating_min: float = 0,
        rating_max: float = 5,
        sort: str = 'default',
        not_empty: str = "false"
):
    '''Страница с результатами поиска'''
    response = product_search_api.predict_for_search_page(
        search_string, limit, price_min, price_max, rating_min, rating_max, sort, not_empty, .2)
    data = response['data']

    return paginate(data)


@router.get('/job_search/', response_model=SearchJob)
async def get_job(job_name: str, limit: int = 10):
    '''Поиск по названию вакансии'''
    answer = job_search_api.predict(job_name, limit, .2)
    return {'data': answer}
