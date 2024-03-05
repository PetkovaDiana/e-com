from .recommendation import Recommendation
from .api_parser import get_products


class RecommendationApi(Recommendation):
    '''Апи интерфейс рекомендаций'''

    def predict(self, product_uuid: str, limit: int) -> list:
        uuids_str: str = ','.join(self.get(product_uuid, limit))
        products = get_products(uuids_str)
        return products