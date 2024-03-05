import type { MarkerClusterer } from "@/types/Cluster.types";
import type {
  KaKaoMapMouseEvent,
  KakaoMap,
  KakaoMarker,
} from "@/types/KakaoMap.types";
import type { Marker } from "@/types/Marker.types";
import AddIcon from "@mui/icons-material/Add";
import GpsOffIcon from "@mui/icons-material/GpsOff";
import LoginIcon from "@mui/icons-material/Login";
import MyLocationIcon from "@mui/icons-material/MyLocation";
import RemoveIcon from "@mui/icons-material/Remove";
import Button from "@mui/material/Button";
import CircularProgress from "@mui/material/CircularProgress";
import { Suspense, lazy, useEffect, useRef, useState } from "react";
import getAllMarker from "../../api/markers/getAllMarker";
import activeMarkerImage from "../../assets/images/cb1.webp";
import pendingMarkerImage from "../../assets/images/cb2.webp";
import useDeleteUser from "../../hooks/mutation/user/useDeleteUser";
import useMap from "../../hooks/useMap";
import useModalStore from "../../store/useModalStore";
import useUploadFormDataStore from "../../store/useUploadFormDataStore";
import useUserStore from "../../store/useUserStore";
import ActionButton from "../ActionButton/ActionButton";
import AddChinupSkeleton from "../AddChinupBarForm/AddChinupSkeleton";
import CenterBox from "../CenterBox/CenterBox";
import FloatingButton from "../FloatingButton/FloatingButton";
import MarkerInfoSkeleton from "../MarkerInfoModal/MarkerInfoSkeleton";
import BasicModal from "../Modal/Modal";
import MyInfoModal from "../MyInfoModal/MyInfoModal";
import * as Styled from "./Map.style";
import useToastStore from "../../store/useToastStore";

const AddChinupBarForm = lazy(
  () => import("../AddChinupBarForm/AddChinupBarForm")
);
const MarkerInfoModal = lazy(
  () => import("../MarkerInfoModal/MarkerInfoModal")
);

export interface MarkerInfo extends Omit<Marker, "photos"> {
  index: number;
  photos?: string[];
}

const Map = () => {
  const modalState = useModalStore();
  const userState = useUserStore();
  const formState = useUploadFormDataStore();
  const toastState = useToastStore();

  const { mutateAsync } = useDeleteUser();

  const mapRef = useRef<HTMLDivElement | null>(null);
  const map = useMap(mapRef);

  const [isMarked, setIsMarked] = useState(false); // 위치 등록하기 토스트 모달 표시 여부
  const [openForm, setOpenForm] = useState(false); // 위치 등록 폼 모달 표시 여부

  const [markers, setMarkers] = useState<KakaoMarker[]>([]); // 실제 화면에 표시되는 마커들

  const [markerInfoModal, setMarkerInfoModal] = useState(false); // 마커 정보 모달 표시 여부

  const [currentMarkerInfo, setCurrentMarkerInfo] = useState<MarkerInfo | null>(
    null
  ); // 마커 클릭시 표시되는 상세 정보 (클릭하면 변경)

  const [marker, setMarker] = useState<KakaoMarker | null>(null); // 위치 등록 예정인 마커

  const [clusterer, setClusterer] = useState<MarkerClusterer | null>(null); // 클러스터 인스턴스

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(false);

  const [myInfoModal, setMyInfoModal] = useState(false); // 내 정보 모달 여부
  const [deleteUserModal, setDeleteUserModal] = useState(false); // 회원 탈퇴 모달 여부
  const [deleteUserLoading, setDeleteUserLoading] = useState(false); // 회원 탈퇴 모달 여부

  const imageSize = new window.kakao.maps.Size(39, 39);
  const imageOption = { offset: new window.kakao.maps.Point(27, 45) };

  useEffect(() => {
    toastState.close();
    toastState.setToastText("");
  }, []);

  useEffect(() => {
    setLoading(true);
    setError(false);

    if (!map) return;

    // 클러스터 인스턴스 생성 및 저장
    const clusterer = new window.kakao.maps.MarkerClusterer({
      map: map,
      averageCenter: true,
      minLevel: 10,
    });
    setClusterer(clusterer);

    const activeMarkerImg = new window.kakao.maps.MarkerImage(
      activeMarkerImage,
      imageSize,
      imageOption
    );

    const pendingMarkerImg = new window.kakao.maps.MarkerImage(
      pendingMarkerImage,
      imageSize,
      imageOption
    );

    // 클릭 위치 마커 생성 및 표시
    const clickMarker = new window.kakao.maps.Marker({
      image: pendingMarkerImg,
    });
    clickMarker.setMap(map);
    // 첫 로딩시 화면에서 숨김
    clickMarker.setVisible(false);
    setMarker(clickMarker);

    getAllMarker()
      .then((res) => {
        setError(false);

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
            const images = markerData.photos?.map((photo) => photo.photoUrl);

            setMarkerInfoModal(true);
            setCurrentMarkerInfo({
              ...markerData,
              photos: images || undefined,
              index: index,
            });
          });

          setMarkers((prev) => {
            const copy = [...prev];
            copy.push(newMarker);
            return copy;
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

    window.kakao.maps.event.addListener(
      map,
      "click",
      (mouseEvent: KaKaoMapMouseEvent) => {
        setIsMarked(true);

        const latlng = mouseEvent.latLng;

        clickMarker.setPosition(latlng);

        clickMarker.setMap(map);
        clickMarker.setVisible(true);

        formState.setPosition(latlng.getLat(), latlng.getLng());
      }
    );
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
  const handleMyInfo = () => {
    setMyInfoModal(true);
  };

  const zoomIn = () => {
    const level = map?.getLevel();

    map?.setLevel((level as number) - 1);
  };

  const zoomOut = () => {
    const level = map?.getLevel();

    map?.setLevel((level as number) + 1);
  };

  const handleDeleteUser = async () => {
    setDeleteUserLoading(true);
    try {
      await mutateAsync();
      userState.resetUser();

      setMyInfoModal(false);
      setDeleteUserModal(false);

      toastState.setToastText("탈퇴 완료");
      toastState.open();
    } catch (error) {
      console.log(error);
    } finally {
      setDeleteUserLoading(false);
    }
  };

  return (
    <div>
      <Styled.MapContainer ref={mapRef} />
      {loading && (
        <CenterBox bg="black">
          <CircularProgress
            size={50}
            sx={{
              color: "#fff",
            }}
          />
        </CenterBox>
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
          <Suspense fallback={<AddChinupSkeleton />}>
            <AddChinupBarForm
              setState={setOpenForm}
              setIsMarked={setIsMarked}
              setMarkerInfoModal={setMarkerInfoModal}
              setCurrentMarkerInfo={setCurrentMarkerInfo}
              setMarkers={setMarkers}
              markers={markers}
              map={map}
              marker={marker}
              clusterer={clusterer as MarkerClusterer}
            />
          </Suspense>
        </BasicModal>
      )}
      {markerInfoModal && (
        <BasicModal setState={setMarkerInfoModal}>
          <Suspense fallback={<MarkerInfoSkeleton />}>
            <MarkerInfoModal
              currentMarkerInfo={currentMarkerInfo as MarkerInfo}
              setMarkerInfoModal={setMarkerInfoModal}
              markers={markers}
              setMarkers={setMarkers}
              clusterer={clusterer as MarkerClusterer}
            />
          </Suspense>
        </BasicModal>
      )}

      {deleteUserModal && (
        <BasicModal setState={setDeleteUserModal}>
          <p>정말 탈퇴하시겠습니까?</p>
          <Styled.DeleteUserButtonsWrap>
            <ActionButton bg="black" onClick={handleDeleteUser}>
              {deleteUserLoading ? (
                <CircularProgress size={20} sx={{ color: "#fff" }} />
              ) : (
                "탈퇴하기"
              )}
            </ActionButton>
            <ActionButton
              bg="gray"
              onClick={() => {
                setDeleteUserModal(false);
              }}
            >
              취소
            </ActionButton>
          </Styled.DeleteUserButtonsWrap>
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
        tooltip={userState.user.token ? "메뉴" : "로그인"}
        onClickFn={
          userState.user.token
            ? myInfoModal
              ? () => {
                  setMyInfoModal(false);
                }
              : handleMyInfo
            : handleOpen
        }
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
      {myInfoModal && (
        <MyInfoModal
          map={map as KakaoMap}
          setMyInfoModal={setMyInfoModal}
          setDeleteUserModal={setDeleteUserModal}
        />
      )}
    </div>
  );
};

export default Map;
