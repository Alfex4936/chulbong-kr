from chulbong_fastapi.db.session import SessionLocal


async def get_db():
    async with SessionLocal() as session:
        yield session
