import pandas as pd
from sqlalchemy import create_engine

import os


class PandasDataBase():
    '''Работа с БД для Pandas'''
    port = os.environ.get('DB_PORT')
    ip = os.environ.get('DB_HOST')
    user_name = os.environ.get('DB_USER')
    password = os.environ.get('DB_PASSWORD')
    db_name = os.environ.get('DB_NAME')

    products_table = 'product'
    jobs_table = 'vacancy'

    def __init__(self):
        print("RecommendationService: PandasDataBase.__init__(): start connection")
        self.connection = create_engine(
            'postgresql+psycopg2://{user_name}:{password}@{ip}:{port}/{db_name}'.format(
                user_name=self.user_name,
                password=self.password,
                ip=self.ip,
                port=self.port,
                db_name=self.db_name
            )
        )
        print("RecommendationService: PandasDataBase.__init__(): created connection")

    def _get(self, table_name: str, index_col: str, columns: list, chunksize: int) -> pd.DataFrame:
        print("RecommendationService: PandasDataBase._get(): start load DataFrame")
        chunks = pd.read_sql_table(
            table_name=table_name,
            con=self.connection,
            index_col=index_col,
            columns=columns,
            chunksize=chunksize
        )
        df = pd.concat(chunks)
        print("RecommendationService: PandasDataBase._get(): loaded DataFrame")
        print(f"RecommendationService: PandasDataBase._get(): df.shape = {df.shape}")
        return df

    def get_products(self) -> pd.DataFrame:
        '''Get products from database'''
        return self._get(self.products_table, 'uuid', ['title', 'vendor_code'], 10_000)
