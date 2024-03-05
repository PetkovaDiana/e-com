import uvicorn
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from .database import engine
from .models import Base
from .routers import router
from .tags import description, tags_metadata

Base.metadata.create_all(bind=engine)

app = FastAPI(
    title="BdAPIApp",
    description=description,
    version="0.2",
    openapi_tags=tags_metadata
)

origins = [
    "http://localhost:3000",
    "https://ufaelectro.ru"
]

app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(router)