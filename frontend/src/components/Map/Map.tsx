import type { KaKaoMapMouseEvent } from "@/types/KakaoMap.types";
import AddIcon from "@mui/icons-material/Add";
import GpsOffIcon from "@mui/icons-material/GpsOff";
import LoginIcon from "@mui/icons-material/Login";
import MyLocationIcon from "@mui/icons-material/MyLocation";
import RemoveIcon from "@mui/icons-material/Remove";
import { Button } from "@mui/material";
import { useEffect, useRef, useState } from "react";
import customMarkerImage from "../../assets/images/cb1.png";
import useMap from "../../hooks/useMap";
import useModalStore from "../../store/useModalStore";
import useUserStore from "../../store/useUserStore";
import AddChinupBarForm from "../AddChinupBarForm/AddChinupBarForm";
import FloatingButton from "../FloatingButton/FloatingButton";
import BasicModal from "../Modal/Modal";
import * as Styled from "./Map.style";

const Map = () => {
  const modalState = useModalStore();
  const userState = useUserStore();

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

  const resetCurrentPosition = () => {
    const moveLatLon = new window.kakao.maps.LatLng(37.566535, 126.9779692);
    map?.setCenter(moveLatLon);
  };

  const handleOpen = () => {
    modalState.openLogin();
  };
  const handleLogout = () => {
    userState.resetUser();
  };

  const zoomIn = () => {
    const level = map?.getLevel();

    map?.setLevel((level as number) - 1);
  };

  const zoomOut = () => {
    const level = map?.getLevel();

    map?.setLevel((level as number) + 1);
  };

  return (
    <div>
      <Styled.MapContainer ref={mapRef} />
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
      <FloatingButton
        text={
          userState.user.token ? (
            userState.user.user.email[0].toUpperCase()
          ) : (
            <LoginIcon />
          )
        }
        top={20}
        right={20}
        shape="circle"
        tooltip={userState.user.token ? "로그아웃" : "로그인"}
        onClickFn={userState.user.token ? handleLogout : handleOpen}
      />

      <FloatingButton
        text={<MyLocationIcon />}
        top={200}
        right={20}
        tooltip="내 위치"
        onClickFn={centerMapOnCurrentPosition}
      />
      <FloatingButton
        text={<GpsOffIcon />}
        top={240}
        right={20}
        tooltip="위치 초기화"
        onClickFn={resetCurrentPosition}
      />

      <FloatingButton
        text={<AddIcon />}
        top={300}
        right={20}
        tooltip="확대"
        onClickFn={zoomIn}
      />
      <FloatingButton
        text={<RemoveIcon />}
        top={340}
        right={20}
        tooltip="축소"
        onClickFn={zoomOut}
      />
    </div>
  );
};

export default Map;
