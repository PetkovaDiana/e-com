from pydantic import BaseModel
from typing import List

from .funcs import beautiful, letter_price


class CartProductAPI(BaseModel):
    '''Продукт в корзине'''
    id: int
    article: str
    title: str
    count: float
    unit: str
    price: float
    cost_price: float

    def get_cost_price(self):
        return beautiful(self.cost_price )


class CustomerAPI(BaseModel):
    '''Покупатель'''
    title: str
    inn: str
    kpp: str
    address: str
    email: str


class CartAPI(BaseModel):
    '''Корзина'''
    cart_products: List[CartProductAPI]
    total_price: float
    nds: float
    product_counter: int

    def get_total_price(self):
        return beautiful(self.total_price)

    def get_nds(self):
        return beautiful(self.nds)

    def get_letter_price(self):
        return letter_price(self.total_price)


class OrderAPI(BaseModel):
    '''Заказ'''
    number: str
    date: str
    customer: CustomerAPI
    cart: CartAPI


class AnswerAPI(BaseModel):
    '''Ответ на отправку счета'''
    status: int
    url: str
