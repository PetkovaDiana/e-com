from sqlalchemy import create_engine
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker

import os


SQLALCHEMY_DATABASE_URL = "postgresql://{login}:{password}@{ip}/{name}".format(
    login=os.environ.get('DB_USER'),
    password=os.environ.get('DB_PASSWORD'),
    ip=os.environ.get('DB_HOST'),
    name=os.environ.get('DB_NAME'),
    post=os.environ.get('DB_PORT')
)

engine = create_engine(
    SQLALCHEMY_DATABASE_URL
)

Base = declarative_base()

SessionLocal = sessionmaker(autoflush=False, bind=engine)


def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()