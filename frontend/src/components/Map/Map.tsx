import type { KaKaoMapMouseEvent, KakaoMarker } from "@/types/KakaoMap.types";
import AddIcon from "@mui/icons-material/Add";
import GpsOffIcon from "@mui/icons-material/GpsOff";
import LoginIcon from "@mui/icons-material/Login";
import MyLocationIcon from "@mui/icons-material/MyLocation";
import RemoveIcon from "@mui/icons-material/Remove";
import { Button, CircularProgress } from "@mui/material";
import { useEffect, useRef, useState, useCallback } from "react";
import getAllMarker from "../../api/markers/getAllMarker";
import activeMarkerImage from "../../assets/images/cb1.png";
import pendingMarkerImage from "../../assets/images/cb2.png";
import useMap from "../../hooks/useMap";
import useModalStore from "../../store/useModalStore";
import useUploadFormDataStore from "../../store/useUploadFormDataStore";
import useUserStore from "../../store/useUserStore";
import AddChinupBarForm from "../AddChinupBarForm/AddChinupBarForm";
import BackgroundBlack from "../BackgroundBlack/BackgroundBlack";
import FloatingButton from "../FloatingButton/FloatingButton";
import MarkerInfoModal from "../MarkerInfoModal/MarkerInfoModal";
import BasicModal from "../Modal/Modal";
import * as Styled from "./Map.style";
import { Marker, Photo } from "@/types/Marker.types";
import { MarkerClusterer } from "@/types/Cluster.types";

export interface MarkerInfo extends Omit<Marker, "photos"> {
  index: number;
  photos?: string[];
}

const Map = () => {
  const modalState = useModalStore();
  const userState = useUserStore();
  const formState = useUploadFormDataStore();

  const mapRef = useRef<HTMLDivElement | null>(null);
  const map = useMap(mapRef);

  const [isMarked, setIsMarked] = useState(false);
  const [openForm, setOpenForm] = useState(false);

  const [markers, setMarkers] = useState<KakaoMarker[]>([]);

  const [markerInfoModal, setMarkerInfoModal] = useState(false);

  const [currentMarkerInfo, setCurrentMarkerInfo] = useState<MarkerInfo | null>(
    null
  );

  const [marker, setMarker] = useState<KakaoMarker | null>(null);

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(false);
  const [clusterer, setClusterer] = useState<MarkerClusterer | null>(null);

  const imageSize = new window.kakao.maps.Size(50, 59);
  const imageOption = { offset: new window.kakao.maps.Point(27, 45) };

  const onClick = useCallback((mouseEvent: KaKaoMapMouseEvent) => {
    if (!map) return;

    const latlng = mouseEvent.latLng;

    // Remove existing marker if there is one
    if (marker) {
      marker.setMap(null);
    }

    const pendingMarkerImg = new window.kakao.maps.MarkerImage(
      pendingMarkerImage,
      imageSize,
      imageOption
    );

    const clickMarker = new window.kakao.maps.Marker({
      position: latlng,
      map: map,
      image: pendingMarkerImg,
    });

    // Update the marker state
    setMarker(clickMarker);
    setIsMarked(true); // Indicate that a marker has been placed

    formState.setPosition(latlng.getLat(), latlng.getLng());
  }, [map, marker, pendingMarkerImage, imageSize, imageOption, setMarker, setIsMarked, formState]);

  // Setting up markers and clusterer
  useEffect(() => {
    if (!map) return;

    setLoading(true);
    setError(false);

    const clusterer = new window.kakao.maps.MarkerClusterer({
      map: map,
      averageCenter: true,
      minLevel: 10,
    });

    setClusterer(clusterer);

    window.kakao.maps.event.addListener(map, "click", onClick);

    // Cleanup function to remove the event listener
    return () => {
      window.kakao.maps.event.removeListener(map, "click", onClick);
    };
  }, [map, marker, setIsMarked, setMarker, formState.setPosition]);

  // Fetching and displaying markers
  useEffect(() => {
    if (!map || !clusterer) return;

    getAllMarker()
      .then((res) => {
        setError(false);

        const activeMarkerImg = new window.kakao.maps.MarkerImage(
          activeMarkerImage,
          imageSize,
          imageOption
        );

        res.forEach((markerData, index) => {
          const markerPosition = new window.kakao.maps.LatLng(
            markerData.latitude,
            markerData.longitude
          );
          const newMarker = new window.kakao.maps.Marker({
            map: map,
            position: markerPosition,
            title: markerData.description,
            image: activeMarkerImg,
          });

          window.kakao.maps.event.addListener(newMarker, "click", () => {
            const images = markerData.photos?.map(photo => photo.photoUrl);
            setMarkerInfoModal(true);
            setCurrentMarkerInfo({
              ...markerData,
              photos: images || [],
              index: index
            });
          });

          clusterer.addMarker(newMarker);
        });

      })
      .catch((error) => {
        setError(true);
        console.error(error);
      })
      .finally(() => {
        setLoading(false);
      });

    // if (clusterer) { 미작동
    //   clusterer.redraw();
    // }
  }, [map, clusterer, getAllMarker]);


  const centerMapOnCurrentPosition = () => {
    if (map && navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          const moveLatLon = new window.kakao.maps.LatLng(
            position.coords.latitude,
            position.coords.longitude
          );
          map.setCenter(moveLatLon);
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
      {loading && (
        <BackgroundBlack>
          <CircularProgress
            size={50}
            sx={{
              color: "#fff",
            }}
          />
        </BackgroundBlack>
      )}

      {!loading && error && (
        <BasicModal>
          <div>철봉 마커를 가져오는 데 실패하였습니다.</div>
          <Button
            onClick={() => {
              window.location.reload();
            }}
            sx={{
              color: "#fff",
              width: "100%",
              height: "40px",
              backgroundColor: "#333",
              marginTop: "1rem",
              "&:hover": {
                backgroundColor: "#555",
              },
            }}
          >
            다시 시도하기
          </Button>
        </BasicModal>
      )}

      {openForm && (
        <BasicModal setState={setOpenForm}>
          <AddChinupBarForm
            setState={setOpenForm}
            setIsMarked={setIsMarked}
            setMarkerInfoModal={setMarkerInfoModal}
            setCurrentMarkerInfo={setCurrentMarkerInfo}
            setMarkers={setMarkers}
            markers={markers}
            map={map}
            marker={marker}
          />
        </BasicModal>
      )}
      {markerInfoModal && (
        <BasicModal setState={setMarkerInfoModal}>
          <MarkerInfoModal
            currentMarkerInfo={currentMarkerInfo as MarkerInfo}
            setMarkerInfoModal={setMarkerInfoModal}
            markers={markers}
          />
        </BasicModal>
      )}
      <Button
        onClick={() => {
          if (userState.user.token === "") {
            modalState.openLogin();
          } else {
            setOpenForm(true);
          }
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
            marker?.setMap(null);
            setMarker(null);
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
