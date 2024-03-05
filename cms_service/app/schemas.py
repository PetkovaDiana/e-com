from typing import Optional
from pydantic import BaseModel
from datetime import datetime


class MainModel(BaseModel):
    id: int
    title: Optional[str]
    link: Optional[str]

    class Config:
        orm_mode = True


class MainPageBannerDTO(MainModel):
    description: Optional[str]
    image_right: Optional[str]


class BlogDTO(MainModel):
    description: Optional[str]
    short_description: Optional[str]
    image: Optional[str]
    date: Optional[datetime]


class PromotionsDTO(MainModel):
    description: Optional[str]
    image: Optional[str]
    date: Optional[datetime]


class Footer(BaseModel):
    id: int

    class Config:
        orm_mode = True


class PhonesDTO(Footer):
    phone: Optional[str]


class HeaderPhonesDTO(Footer):
    phone: Optional[str]


class AddressesDTO(Footer):
    address: Optional[str]


class ObjectsDTO(Footer):
    icon: Optional[str]
    link: Optional[str]


class MetaTagsDTO(Footer):
    '''Мета тэги для index.html'''
    title: Optional[str]
    tag: Optional[str]
