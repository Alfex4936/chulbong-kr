import MyLocationIcon from "@mui/icons-material/MyLocation";
import Fab from "@mui/material/Fab";
import { useEffect, useRef } from "react";
import useMap from "../../hooks/useMap";
import * as Styled from "./Map.style";

const Map = () => {
  const mapRef = useRef();
  const map = useMap(mapRef);

  useEffect(() => {
    if (map) {
      kakao.maps.event.addListener(map, "click", (mouseEvent) => {
        const latlng = mouseEvent.latLng;

        const marker = new kakao.maps.Marker({
          position: map.getCenter(),
        });

        marker.setMap(map);

        marker.setPosition(latlng);

        let message = "클릭한 위치의 위도는 " + latlng.getLat() + " 이고, ";
        message += "경도는 " + latlng.getLng() + " 입니다";

        console.log(message);
      });
    }
  }, [map]);

  const centerMapOnCurrentPosition = () => {
    if (map && navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          const newPos = new window.kakao.maps.LatLng(
            position.coords.latitude,
            position.coords.longitude
          );
          map.setLevel(3);
          map.setCenter(newPos);

          console.log(position);
        },
        (error) => {
          console.error(error);
        }
      );
    } else {
      alert(
        "Geolocation is not supported by this browser or map is not loaded yet."
      );
    }
  };

  return (
    <div>
      <Styled.MapContainer ref={mapRef} />
      <Fab
        color="secondary"
        aria-label="locate"
        onClick={() => {
          centerMapOnCurrentPosition();
        }}
        sx={{
          position: "absolute",
          bottom: 32,
          right: 32,
          color: "white",
          bgcolor: "black",
          "&:hover": {
            bgcolor: "gray",
          },
          boxShadow: "0px 0px 10px rgba(0, 0, 0, 0.5)",
          border: "2px solid white",
        }}
      >
        <MyLocationIcon />
      </Fab>
    </div>
  );
};

export default Map;