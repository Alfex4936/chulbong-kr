from contextlib import asynccontextmanager

from fastapi import FastAPI

from chulbong_fastapi.db.session import database


@asynccontextmanager
async def lifespan(_: FastAPI):
    await database.connect()
    yield
    await database.disconnect()
