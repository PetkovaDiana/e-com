from .api_parser import get_products
from .utils import clean_str
from .job_search import JobsSearch
from .product_search import ProductSearch


class SearchProductsApi(ProductSearch):
    '''Апи интерфейс поисковика товаров'''

    def predict(self, string_search: str, size: int, cnt_typos: float, to_products_page=False) -> list:
        '''Поиск
        string_search - строка для поиска
        size - количество ответов
        cnt_typos - доля опечаток
        '''

        # string_search = clean_str(string_search)
        string_search = string_search.lower() # TODO fix me
        string_search = string_search.replace('иек', 'iek')
        json_answer, count = [], 0
        ans = set()

        # Автодополнение по артикулу
        _products = self.vendor_code_trie.predict(string_search)
        if _products:
            for product in _products:
                if product.title:
                    if product.uuid in ans:
                        continue
                    else:
                        if to_products_page:
                            json_answer.append(product.uuid)
                        else:
                            json_answer.append({
                                'id': str(product.uuid),
                                'title': str(product.original_title),
                                'vendor_code': str(product.title)
                            })
                            count += 1
                            if count > size:
                                break

                        ans.add(product.uuid)

        # Автодополнение по полному названию
        _products = self.products_trie.predict(string_search)
        if _products:
            for product in _products:
                if product.uuid in ans:
                    continue
                else:
                    if to_products_page:
                        json_answer.append(product.uuid)
                    else:
                        json_answer.append({
                            'id': product.uuid,
                            'title': product.original_title,
                            'vendor_code': product.vendor_code
                        })
                        count += 1
                        if count > size:
                            break
                    ans.add(product.uuid)

        indexes = self._words_search(string_search, cnt_typos)
        indexes = sorted(indexes.items(), key=lambda item: item[-1], reverse=True)

        for index in indexes:
            product = self.products.loc[index[0]]
            if index[0] in ans:
                continue
            else:
                if to_products_page:
                    json_answer.append(index[0])
                else:
                    json_answer.append({
                        'id': index[0],
                        'title': product['original_title'],
                        'vendor_code': product['vendor_code']
                    })
                    count += 1
                    if count > size:
                        break
                ans.add(index[0])

        return json_answer[:size]

    def predict_for_search_page(self, search_string: str, size: int, price_min: float,
                                price_max: float, rating_min: float, rating_max: float,
                                sort: str, not_empty: str, cnt_typos: float) -> dict:
        uuids_list = self.predict(search_string, size, cnt_typos, True)
        uuids: str = ','.join(uuids_list)
        products = get_products(uuids, price_min, price_max, rating_min, rating_max, sort, not_empty)
        return products


class SearchJobsApi(JobsSearch):
    '''Апи интерфейс поисковика вакансий'''

    def predict(self, search_string: str, size: int, cnt_typos: float) -> list:
        string_search = clean_str(search_string)
        json_answer, count = [], 0

        _jobs = self.jobs_trie.predict(string_search)
        if _jobs:
            for job in _jobs:
                json_answer.append({
                    'id': job.id,
                    'title': job.title,
                    'first_phone': job.first_phone if job.first_phone else '',
                    'second_phone': job.second_phone if job.second_phone else '',
                    'email': job.email if job.email else ''
                })
                count += 1
                if count == size:
                    break

        indexes = self._words_search(string_search, cnt_typos)
        indexes = sorted(indexes.items(), key=lambda item: item[-1], reverse=True)

        for index in indexes:
            job = self.jobs.loc[index[0]]
            json_answer.append({
                'id': int(index[0]),
                'title': job.original_title,
                'first_phone': job['first_phone'] if job['first_phone'] else '',
                'second_phone': job['second_phone'] if job['second_phone'] else '',
                'email': job['email'] if job['email'] else ''
            })
            count += 1
            if count == size:
                break

        return json_answer