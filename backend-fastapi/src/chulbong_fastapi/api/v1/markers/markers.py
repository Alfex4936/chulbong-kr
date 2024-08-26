from typing import List

from fastapi import APIRouter, Depends
from geoalchemy2.shape import to_shape
from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession

from chulbong_fastapi.api.deps import get_db
from chulbong_fastapi.db.models import Marker
from chulbong_fastapi.schemas.marker import MarkerSchema

router = APIRouter()


@router.get("/test")
async def test_api() -> str:
    return "hello"


@router.get("/", response_model=List[MarkerSchema])
async def read_markers(db: AsyncSession = Depends(get_db)):
    result = await db.execute(select(Marker))
    markers = result.scalars().all()
    return [
        MarkerSchema(
            MarkerID=marker.MarkerID,
            UserID=marker.UserID,
            Description=marker.Description,
            CreatedAt=marker.CreatedAt,
            UpdatedAt=marker.UpdatedAt,
            Address=marker.Address,
            Location=to_shape(marker.Location).wkt if marker.Location else None,
        )
        for marker in markers
    ]
