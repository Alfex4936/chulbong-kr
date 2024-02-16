import type { KaKaoMapMouseEvent } from "@/types/KakaoMap.types";
import MyLocationIcon from "@mui/icons-material/MyLocation";
import { Button } from "@mui/material";
import Fab from "@mui/material/Fab";
import { useEffect, useRef, useState } from "react";
import customMarkerImage from "../../assets/images/cb1.png";
import useMap from "../../hooks/useMap";
import AddChinupBarForm from "../AddChinupBarForm/AddChinupBarForm";
import BasicModal from "../Modal/Modal";
import * as Styled from "./Map.style";

const Map = () => {
  const mapRef = useRef<HTMLDivElement | null>(null);
  const map = useMap(mapRef);

  const [isMarked, setIsMarked] = useState(false);
  const [openForm, setOpenForm] = useState(false);

  useEffect(() => {
    if (map) {
      const imageSize = new window.kakao.maps.Size(50, 59);
      const imageOption = { offset: new window.kakao.maps.Point(27, 45) };

      const markerImage = new window.kakao.maps.MarkerImage(
        customMarkerImage,
        imageSize,
        imageOption
      );

      const marker = new window.kakao.maps.Marker({
        image: markerImage,
      });

      marker.setMap(map);

      window.kakao.maps.event.addListener(
        map,
        "click",
        (mouseEvent: KaKaoMapMouseEvent) => {
          setIsMarked(true);

          const latlng = mouseEvent.latLng;

          marker.setPosition(latlng);

          let message = `클릭한 위치의 위도는 ${latlng.getLat()} 이고, `;
          message += `경도는 ${latlng.getLng()} 입니다`;

          console.log(message);
        }
      );
    }
  }, [map]);

  const centerMapOnCurrentPosition = () => {
    if (map && navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          const moveLatLon = new window.kakao.maps.LatLng(
            position.coords.latitude,
            position.coords.longitude
          );
          map.setCenter(moveLatLon);
          console.log(position.coords.latitude, position.coords.longitude);
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
        onClick={centerMapOnCurrentPosition}
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
      {openForm && (
        <BasicModal setState={setOpenForm}>
          <AddChinupBarForm />
        </BasicModal>
      )}
      <Button
        onClick={() => {
          setOpenForm(true);
        }}
        sx={{
          position: "absolute",
          opacity: map && isMarked ? "100" : "0",
          bottom: "30px",
          left: map && isMarked ? "50%" : "10%",
          transform: "translateX(-50%)",
          transition: "all .3s",
          color: "#fff",
          backgroundColor: "#333",
          zIndex: "1",
          width: "300px",
          height: "60px",
          "&:hover": {
            backgroundColor: "#555",
          },
        }}
      >
        <Styled.ExitButton
          onClick={(e) => {
            e.stopPropagation();
            setIsMarked(false);
          }}
        >
          X
        </Styled.ExitButton>
        위치 등록하기
      </Button>
    </div>
  );
};

export default Map;
