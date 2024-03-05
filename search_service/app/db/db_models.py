from typing import List

from pydantic import BaseModel


class SearchWindowProduct(BaseModel):
    '''Товар для модального окна с авто дополнением поиска'''
    id: str
    title: str
    vendor_code: str


class SearchWindow(BaseModel):
    data: List[SearchWindowProduct]


class CountProducts(BaseModel):
    '''Количество товара в запросе'''
    count: str


class Products(BaseModel):
    '''Товары на странице поисковой выдачи'''
    id: str
    title: str
    base_unit: str
    count: str
    image: str
    price: str
    quantity: str
    rating: str
    review_count: str
    vendor_code: str


class Job(BaseModel):
    '''Модель вакансии'''
    id: int
    title: str
    first_phone: str
    second_phone: str
    email: str


class SearchJob(BaseModel):
    data: List[Job]