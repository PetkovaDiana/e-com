from requests import get

URL = 'http://127.0.0.1:8000/search/'

def search(search_string: str, limit: int) -> dict:
    params = {
        'search_string': search_string,
        'limit': limit
    }
    response = get(URL, params=params)
    assert response.status_code == 200, "Ошибка в запросе"
    return response.json()['data'][0]


def test_vendor_codes():
    vendor_codes = ['MVA20-1-016-C', 'sq1806-0112', 'SQ0224-0025', '14210dek']
    for vendor_code in vendor_codes:
        search_result = search(vendor_code, 1)['vendor_code']
        assert vendor_code.lower() == search_result, 'Неверный результат поискового запроса'


def test_titles_1():
    '''Поиск подстроки в строке имеющей спецсимволы'''
    titles = ['3x1.5']
    for title in titles:
        search_result = search(title, 1)['title']
        assert title.lower() in search_result, 'Неверный результат поискового запроса'


def test_titles_2():
    '''Поиск по полному названию с полным совпадением'''
    titles = ['кабель кгввнг(а)-ls 3x1', 'кабель силовой кгввнг(а)-ls 3х1,5- 0,66 тртс']
    for title in titles:
        search_result = search(title, 1)['title']
        assert title.lower() == search_result, 'Неверный результат поискового запроса'


if __name__ == "__main__":
    test_vendor_codes()
    test_titles_1()
    test_titles_2()