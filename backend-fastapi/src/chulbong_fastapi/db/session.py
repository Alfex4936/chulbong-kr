import os

from databases import Database
from sqlalchemy.ext.asyncio import AsyncSession, create_async_engine
from sqlalchemy.orm import sessionmaker

from chulbong_fastapi.core.config import DATABASE_URL, DEBUG

database = Database(DATABASE_URL)
engine = create_async_engine(
    DATABASE_URL,
    echo=DEBUG,
    pool_size=10,  # Maximum number of connections in the pool
    max_overflow=5,  # Maximum number of connections above pool_size
    pool_timeout=30,  # Time to wait for a connection before throwing an error
)
SessionLocal = sessionmaker(
    autocommit=False, autoflush=False, bind=engine, class_=AsyncSession
)
