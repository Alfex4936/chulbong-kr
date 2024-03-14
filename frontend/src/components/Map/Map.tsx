import type { MarkerClusterer } from "@/types/Cluster.types";
import { CustomOverlay } from "@/types/CustomOverlay.types";
import type {
  KaKaoMapMouseEvent,
  KakaoMap,
  KakaoMarker,
} from "@/types/KakaoMap.types";
import AddIcon from "@mui/icons-material/Add";
import GpsOffIcon from "@mui/icons-material/GpsOff";
import LoginIcon from "@mui/icons-material/Login";
import MyLocationIcon from "@mui/icons-material/MyLocation";
import QuestionMarkIcon from "@mui/icons-material/QuestionMark";
import RemoveIcon from "@mui/icons-material/Remove";
import Button from "@mui/material/Button";
import CircularProgress from "@mui/material/CircularProgress";
import { bouncy } from "ldrs";
import { Suspense, lazy, useEffect, useRef, useState } from "react";
import { createRoot } from "react-dom/client";
import activeMarkerImage from "../../assets/images/cb1.webp";
import pendingMarkerImage from "../../assets/images/cb2.webp";
import selectedMarkerImage from "../../assets/images/cb3.webp";
import useRequestPasswordReset from "../../hooks/mutation/auth/useRequestPasswordReset";
import useDeleteUser from "../../hooks/mutation/user/useDeleteUser";
import useGetAllMarker from "../../hooks/query/marker/useGetAllMarker";
import useGetMyInfo from "../../hooks/query/user/useGetMyInfo";
import useInput from "../../hooks/useInput";
import useMap from "../../hooks/useMap";
import useMapPositionStore from "../../store/useMapPositionStore";
import useModalStore from "../../store/useModalStore";
import useOnBoardingStore from "../../store/useOnBoardingStore";
import useToastStore from "../../store/useToastStore";
import useUploadFormDataStore from "../../store/useUploadFormDataStore";
import useUserStore from "../../store/useUserStore";
import ActionButton from "../ActionButton/ActionButton";
import AddChinupSkeleton from "../AddChinupBarForm/AddChinupSkeleton";
import CenterBox from "../CenterBox/CenterBox";
import FloatingButton from "../FloatingButton/FloatingButton";
import Input from "../Input/Input";
import MarkerInfoSkeleton from "../MarkerInfoModal/MarkerInfoSkeleton";
import BasicModal from "../Modal/Modal";
import MyInfoModal from "../MyInfoModal/MyInfoModal";
import * as Styled from "./Map.style";
import MapHeader from "./MapHeader";

import "ldrs/ring";
import OnBoarding from "../OnBoarding/OnBoarding";

const AddChinupBarForm = lazy(
  () => import("../AddChinupBarForm/AddChinupBarForm")
);
const MarkerInfoModal = lazy(
  () => import("../MarkerInfoModal/MarkerInfoModal")
);

export interface MarkerInfo {
  markerId: number;
}

const Map = () => {
  const modalState = useModalStore();
  const userState = useUserStore();
  const formState = useUploadFormDataStore();
  const toastState = useToastStore();
  const mapPosition = useMapPositionStore();
  const onBoardingState = useOnBoardingStore();

  const query = new URLSearchParams(location.search);
  const sharedMarker = query.get("d");
  const sharedMarkerLat = query.get("la");
  const sharedMarkerLng = query.get("lo");

  const emailInput = useInput("");

  const { data, isLoading, isError } = useGetAllMarker();

  const { mutateAsync: deleteUser } = useDeleteUser();
  const { mutateAsync: sendPasswordReset } = useRequestPasswordReset(
    emailInput.value
  );
  const { data: myInfo } = useGetMyInfo();

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

  const [myInfoModal, setMyInfoModal] = useState(false); // 내 정보 모달 여부

  const [deleteUserModal, setDeleteUserModal] = useState(false); // 회원 탈퇴 모달 여부
  const [deleteUserLoading, setDeleteUserLoading] = useState(false); // 회원 탈퇴 로딩

  const [sendOk, setSendOk] = useState(false); // 비밀번호 변경 이메일 전송 여부
  const [changePasswordLoading, setChangePasswordLoading] = useState(false); // 비밀 번호 변경 로딩
  const [emailError, setEmailError] = useState(""); // 비밀번호 변경 에러

  const [currentOverlay, setCurrentOverlay] = useState<CustomOverlay | null>(
    null
  ); // GPS 현재 위치
  const [gpsLoading, setGpsLoading] = useState(false); // GPS 현재 위치

  const imageSize = new window.kakao.maps.Size(39, 39);
  const imageOption = { offset: new window.kakao.maps.Point(27, 45) };

  useEffect(() => {
    if (!onBoardingState.isOnBoarding) return;

    if (onBoardingState.step === 2) {
      setIsMarked(true);
    } else {
      setIsMarked(false);
    }
  }, [onBoardingState.step]);

  useEffect(() => {
    const filtering = (markerId: number) => {
      const imageSize = new window.kakao.maps.Size(39, 39);
      const imageOption = { offset: new window.kakao.maps.Point(27, 45) };

      const selectedMarkerImg = new window.kakao.maps.MarkerImage(
        selectedMarkerImage,
        imageSize,
        imageOption
      );
      const activeMarkerImg = new window.kakao.maps.MarkerImage(
        activeMarkerImage,
        imageSize,
        imageOption
      );

      const marker = markers.find((value) => Number(value.Gb) === markerId);

      markers.forEach((marker) => {
        marker?.setImage(activeMarkerImg);
      });

      marker?.setImage(selectedMarkerImg);
    };

    if (sharedMarker && sharedMarkerLat && sharedMarkerLng) {
      if (sharedMarker) filtering(Number(sharedMarker));
      setMarkerInfoModal(true);
      setCurrentMarkerInfo({ markerId: Number(sharedMarker) });

      const moveLatLon = new window.kakao.maps.LatLng(
        Number(sharedMarkerLat),
        Number(sharedMarkerLng)
      );
      mapPosition.setPosition(Number(sharedMarkerLat), Number(sharedMarkerLng));
      map?.setCenter(moveLatLon);
    }
  }, [map, sharedMarker, sharedMarkerLat, sharedMarkerLng, markers]);

  useEffect(() => {
    const handleKeyDownClose = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        setIsMarked(false);
        marker?.setMap(null);
      }
    };

    window.addEventListener("keydown", handleKeyDownClose);

    toastState.close();
    toastState.setToastText("");

    return () => {
      window.removeEventListener("keydown", handleKeyDownClose);
    };
  }, [marker]);

  useEffect(() => {
    if (!map || !data) return;

    // 클러스터 인스턴스 생성 및 저장
    const clusterer = new window.kakao.maps.MarkerClusterer({
      map: map,
      averageCenter: true,
      minLevel: 6,
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
      zIndex: 4,
    });
    clickMarker.setMap(map);
    // 첫 로딩시 화면에서 숨김
    clickMarker.setVisible(false);

    setMarker(clickMarker);

    const newMarkers = data?.map((marker) => {
      const newMarker = new window.kakao.maps.Marker({
        position: new window.kakao.maps.LatLng(
          marker.latitude,
          marker.longitude
        ),
        image: activeMarkerImg,
        title: marker.markerId,
        zIndex: 4,
      });

      window.kakao.maps.event.addListener(newMarker, "click", () => {
        setMarkerInfoModal(true);
        setCurrentMarkerInfo({
          markerId: marker.markerId,
        });
      });

      return newMarker;
    });

    clusterer.addMarkers(newMarkers);
    setMarkers([...newMarkers]);

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
  }, [map, data]);

  const centerMapOnCurrentPosition = () => {
    if (onBoardingState.step === 4) return;
    if (map && navigator.geolocation) {
      setGpsLoading(true);
      navigator.geolocation.getCurrentPosition(
        (position) => {
          const moveLatLon = new window.kakao.maps.LatLng(
            position.coords.latitude,
            position.coords.longitude
          );

          // Remove the last overlay if it exists
          if (currentOverlay) {
            currentOverlay.setMap(null);
          }

          // Create a div for the marker React component
          const overlayDiv = document.createElement("div");
          const root = createRoot(overlayDiv);
          root.render(<Styled.UserLocationMarker />);

          // Create a custom overlay
          const customOverlay = new window.kakao.maps.CustomOverlay({
            position: moveLatLon,
            content: overlayDiv,
            zIndex: 3,
          });

          customOverlay.setMap(map);
          setCurrentOverlay(customOverlay);

          mapPosition.setPosition(
            position.coords.latitude,
            position.coords.longitude
          );
          map.setCenter(moveLatLon);
          setGpsLoading(false);
        },
        (error) => {
          console.error(error);
          setGpsLoading(false);
        }
      );
    } else {
      alert(
        "Geolocation is not supported by this browser or map is not loaded yet."
      );
    }
  };

  const resetCurrentPosition = () => {
    if (onBoardingState.step === 5) return;
    const moveLatLon = new window.kakao.maps.LatLng(37.566535, 126.9779692);
    mapPosition.setPosition(37.566535, 126.9779692);
    map?.setCenter(moveLatLon);
  };

  const handleOpen = () => {
    if (onBoardingState.step === 3) return;
    modalState.openLogin();
  };
  const handleMyInfo = () => {
    if (onBoardingState.step === 3) return;
    setMyInfoModal(true);
  };

  const zoomIn = () => {
    if (onBoardingState.step === 6) return;
    const level = map?.getLevel();

    map?.setLevel((level as number) - 1);
  };

  const zoomOut = () => {
    if (onBoardingState.step === 7) return;
    const level = map?.getLevel();

    map?.setLevel((level as number) + 1);
  };

  const handleDeleteUser = async () => {
    setDeleteUserLoading(true);
    try {
      await deleteUser();
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

  const handleChangePassword = async () => {
    setChangePasswordLoading(true);
    try {
      await sendPasswordReset();
      setSendOk(true);
    } catch (error) {
      setEmailError("잠시 후 다시 시도해 주세요!");
      console.log(error);
    } finally {
      setChangePasswordLoading(false);
    }
  };

  const handleOpenAddMarkerToast = () => {
    if (onBoardingState.isOnBoarding) return;

    if (!myInfo) {
      modalState.openLogin();
    } else {
      const moveLatLon = new window.kakao.maps.LatLng(
        formState.latitude,
        formState.longitude
      );
      mapPosition.setPosition(formState.latitude, formState.longitude);
      map?.setCenter(moveLatLon);
      setOpenForm(true);
    }
  };

  bouncy.register();

  return (
    <Styled.Container>
      <MapHeader map={map} markers={markers} />
      <Styled.MapContainer ref={mapRef} />
      {(isLoading || gpsLoading) && (
        <CenterBox bg="black">
          <CenterBox bg="black">
            <l-bouncy size="80" speed="1.75" color="white" />
          </CenterBox>
        </CenterBox>
      )}

      {!isLoading && onBoardingState.isOnBoarding && <OnBoarding />}

      {!isLoading && isError && (
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
              map={map}
              marker={marker}
              markers={markers}
              clusterer={clusterer as MarkerClusterer}
            />
          </Suspense>
        </BasicModal>
      )}
      {markerInfoModal && (
        <BasicModal setState={setMarkerInfoModal} exit={false}>
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
          <Styled.AlertText>
            <p>정말 탈퇴하시겠습니까?</p>
            <p>
              추가하신 마커는 유지되고, 작성한 댓글 밑 사진은 모두 삭제됩니다!
            </p>
          </Styled.AlertText>
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

      {modalState.passwordModal && (
        <BasicModal>
          {sendOk ? (
            <p
              style={{
                margin: "1rem",
                fontSize: "1.5rem",
              }}
            >
              이메일을 확인해 주세요!
            </p>
          ) : (
            <>
              <p
                style={{
                  margin: "1rem",
                  fontSize: "1.5rem",
                }}
              >
                비밀번호 변경
              </p>
              <Input
                type="email"
                id="email"
                placeholder="이메일"
                value={emailInput.value}
                onChange={(e) => {
                  emailInput.onChange(e);
                  setEmailError("");
                }}
              />
              <Styled.ErrorBox>{emailError}</Styled.ErrorBox>
              <Styled.ChangePasswordButtonsWrap>
                <ActionButton bg="black" onClick={handleChangePassword}>
                  {changePasswordLoading ? (
                    <CircularProgress size={20} sx={{ color: "#fff" }} />
                  ) : (
                    "메일 보내기"
                  )}
                </ActionButton>
                <ActionButton
                  bg="gray"
                  onClick={() => {
                    modalState.close();
                  }}
                >
                  취소
                </ActionButton>
              </Styled.ChangePasswordButtonsWrap>
            </>
          )}
        </BasicModal>
      )}

      <Button
        onClick={handleOpenAddMarkerToast}
        sx={{
          position: "absolute",
          opacity: map && isMarked ? "100" : "0",
          bottom: "30px",
          left: map && isMarked ? "50%" : "10%",
          transform: "translateX(-50%)",
          transition: "all .3s",
          color: "#fff",
          backgroundColor: "#333",
          zIndex: onBoardingState.step === 2 ? "2000" : "1",
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
      <Styled.LoginButtonWrap
        style={{
          zIndex: onBoardingState.step === 3 ? "10000" : "10",
        }}
      >
        <FloatingButton
          text={myInfo ? myInfo.email[0].toUpperCase() : <LoginIcon />}
          top={0}
          left={0}
          shape="circle"
          tooltip={myInfo ? "메뉴" : "로그인"}
          onClickFn={
            myInfo
              ? myInfoModal
                ? () => {
                    if (onBoardingState.step === 3) return;
                    setMyInfoModal(false);
                  }
                : handleMyInfo
              : handleOpen
          }
        />
      </Styled.LoginButtonWrap>

      <FloatingButton
        text={<MyLocationIcon />}
        top={200}
        right={20}
        zIndex={onBoardingState.step === 4 ? 10000 : 10}
        tooltip="내 위치"
        onClickFn={centerMapOnCurrentPosition}
      />
      <FloatingButton
        text={<GpsOffIcon />}
        top={240}
        right={20}
        zIndex={onBoardingState.step === 5 ? 10000 : 10}
        tooltip="위치 초기화"
        onClickFn={resetCurrentPosition}
      />

      <FloatingButton
        text={<AddIcon />}
        top={300}
        right={20}
        zIndex={onBoardingState.step === 6 ? 10000 : 10}
        tooltip="확대"
        onClickFn={zoomIn}
      />
      <FloatingButton
        text={<RemoveIcon />}
        top={340}
        right={20}
        zIndex={onBoardingState.step === 7 ? 10000 : 10}
        tooltip="축소"
        onClickFn={zoomOut}
      />
      <FloatingButton
        text={<QuestionMarkIcon />}
        shape="circle"
        bottom={250}
        right={20}
        zIndex={onBoardingState.step === 12 ? 10000 : 10}
        tooltip="도움말"
        onClickFn={() => {
          onBoardingState.open();
        }}
      />
      {myInfoModal && (
        <MyInfoModal
          map={map as KakaoMap}
          markers={markers}
          setMyInfoModal={setMyInfoModal}
          setDeleteUserModal={setDeleteUserModal}
        />
      )}
    </Styled.Container>
  );
};

export default Map;
