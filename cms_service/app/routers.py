from fastapi import APIRouter, Depends
from sqlalchemy.orm import Session

from .crud import get_all, get_by_id
from .database import get_db
from .models import MainPageBanner, Blog, Phones, HeaderPhones, Addresses, Objects, Promotions, MetaTags
from .schemas import MainPageBannerDTO, BlogDTO, HeaderPhonesDTO, PhonesDTO, AddressesDTO, ObjectsDTO, PromotionsDTO, \
    MetaTagsDTO

router = APIRouter()

get_db()
path = "/api"


@router.get(path + "/banner", tags=["Get all"], response_model=list[MainPageBannerDTO])
def get_all_banner(db: Session = Depends(get_db)):
    banners = get_all(db, MainPageBanner)

    for banner in banners:
        banner.media_root()

    return get_all(db, MainPageBanner)


@router.get(path + "/meta_tags", tags=["Get all"], response_model=list[MetaTagsDTO])
def get_all_meta_tags(db: Session = Depends(get_db)):
    return get_all(db, MetaTags)


@router.get(path + "/blog", tags=["Get all"], response_model=list[BlogDTO])
def get_all_blogs(db: Session = Depends(get_db)):
    blogs = get_all(db, Blog)

    for blog in blogs:
        blog.media_root()

    return blogs


@router.get(path + "/footer_phones", tags=["Get all"], response_model=list[PhonesDTO])
def get_all_phones(db: Session = Depends(get_db)):
    return get_all(db, Phones)


@router.get(path + "/header_phones", tags=["Get all"], response_model=list[HeaderPhonesDTO])
def get_all_phones(db: Session = Depends(get_db)):
    return get_all(db, HeaderPhones)


@router.get(path + "/footer_addresses", tags=["Get all"], response_model=list[AddressesDTO])
def get_all_addresses(db: Session = Depends(get_db)):
    return get_all(db, Addresses)


@router.get(path + "/footer_objects", tags=["Get all"], response_model=list[ObjectsDTO])
def get_all_objects(db: Session = Depends(get_db)):
    objs = get_all(db, Objects)

    for obj in objs:
        obj.media_root()

    return objs


@router.get(path + "/promotions", tags=["Get all"], response_model=list[PromotionsDTO])
def get_all_promotions(db: Session = Depends(get_db)):
    proms = get_all(db, Promotions)

    for prom in proms:
        prom.media_root()

    return proms


@router.get(path + "/banner/{id}", tags=["Get by id"], response_model=MainPageBannerDTO)
def get_id_banner(id: int, db: Session = Depends(get_db)):
    banner = get_by_id(db, MainPageBanner, id)

    banner.media_root()

    return banner


@router.get(path + "/blog/{id}", tags=["Get by id"], response_model=BlogDTO)
def get_id_blog(id: int, db: Session = Depends(get_db)):
    blog = get_by_id(db, Blog, id)

    blog.media_root()

    return blog


@router.get(path + "/footer_phones/{id}", tags=["Get by id"], response_model=PhonesDTO)
def get_id_phones(id: int, db: Session = Depends(get_db)):
    return get_by_id(db, Phones, id)


@router.get(path + "/header_phones/{id}", tags=["Get by id"], response_model=HeaderPhonesDTO)
def get_id_phones(id: int, db: Session = Depends(get_db)):
    return get_by_id(db, HeaderPhones, id)


@router.get(path + "/footer_addresses/{id}", tags=["Get by id"], response_model=AddressesDTO)
def get_id_addresses(id: int, db: Session = Depends(get_db)):
    return get_by_id(db, Addresses, id)


@router.get(path + "/footer_objects/{id}", tags=["Get by id"], response_model=ObjectsDTO)
def get_id_objects(id: int, db: Session = Depends(get_db)):
    obj = get_by_id(db, Objects, id)

    obj.media_root()

    return obj


@router.get(path + "/promotions/{id}", tags=["Get by id"], response_model=PromotionsDTO)
def get_id_promotions(id: int, db: Session = Depends(get_db)):
    prom = get_by_id(db, Promotions, id)
    prom.media_root()

    return prom
