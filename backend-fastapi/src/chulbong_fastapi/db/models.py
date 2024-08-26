from geoalchemy2 import Geometry
from sqlalchemy import TIMESTAMP, Column, ForeignKey, Integer, String

from chulbong_fastapi.db.base import Base


class Marker(Base):
    __tablename__ = "Markers"
    MarkerID = Column(Integer, primary_key=True, index=True)
    UserID = Column(Integer, ForeignKey("Users.UserID"), nullable=True)
    Location = Column(Geometry("POINT", srid=4326), nullable=False)
    Description = Column(String(255), nullable=True)
    CreatedAt = Column(TIMESTAMP, default="CURRENT_TIMESTAMP")
    UpdatedAt = Column(
        TIMESTAMP, default="CURRENT_TIMESTAMP", onupdate="CURRENT_TIMESTAMP"
    )
    Address = Column(String(255), nullable=True)
