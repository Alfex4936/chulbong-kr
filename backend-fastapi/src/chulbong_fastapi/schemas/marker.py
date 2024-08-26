from datetime import datetime
from typing import Optional

from pydantic import BaseModel


# Same as DTO
class MarkerSchema(BaseModel):
    MarkerID: int
    UserID: Optional[int]
    Description: Optional[str]
    CreatedAt: datetime
    UpdatedAt: datetime
    Address: Optional[str]
    Location: str  # Store the Location as a WKT string or as a dictionary

    class Config:
        orm_mode = True
