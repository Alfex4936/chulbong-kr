import type { MarkerClusterer } from "@/types/Cluster.types";
import type {
  KaKaoMapMouseEvent,
  KakaoMap,
  KakaoMarker,
} from "@/types/KakaoMap.types";
import AddIcon from "@mui/icons-material/Add";
import GpsOffIcon from "@mui/icons-material/GpsOff";
import LoginIcon from "@mui/icons-material/Login";
import MyLocationIcon from "@mui/icons-material/MyLocation";
import RemoveIcon from "@mui/icons-material/Remove";
import Button from "@mui/material/Button";
import CircularProgress from "@mui/material/CircularProgress";
import { bouncy } from "ldrs";
import { Suspense, lazy, useEffect, useRef, useState } from "react";
import activeMarkerImage from "../../assets/images/cb1.webp";
import pendingMarkerImage from "../../assets/images/cb2.webp";
import useRequestPasswordReset from "../../hooks/mutation/auth/useRequestPasswordReset";
import useDeleteUser from "../../hooks/mutation/user/useDeleteUser";
import useGetAllMarker from "../../hooks/query/marker/useGetAllMarker";
import useInput from "../../hooks/useInput";
import useMap from "../../hooks/useMap";
import useMapPositionStore from "../../store/useMapPositionStore";
import useModalStore from "../../store/useModalStore";
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
import * as Styled from "./Map.style";

import "ldrs/ring";

const AddChinupBarForm = lazy(
  () => import("../AddChinupBarForm/AddChinupBarForm")
);
const MarkerInfoModal = lazy(
  () => import("../MarkerInfoModal/MarkerInfoModal")
);

import MyInfoModal from "../MyInfoModal/MyInfoModal";
import MapHeader from "./MapHeader";

export interface MarkerInfo {
  markerId: number;
}

const Map = () => {
  const modalState = useModalStore();
  const userState = useUserStore();
  const formState = useUploadFormDataStore();
  const toastState = useToastStore();
  const mapPosition = useMapPositionStore();

  const emailInput = useInput("");

  const { data, isLoading, isError } = useGetAllMarker();

  const { mutateAsync: deleteUser } = useDeleteUser();
  const { mutateAsync: sendPasswordReset } = useRequestPasswordReset(
    emailInput.value
  );

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

  const imageSize = new window.kakao.maps.Size(39, 39);
  const imageOption = { offset: new window.kakao.maps.Point(27, 45) };

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
    if (map && navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          const moveLatLon = new window.kakao.maps.LatLng(
            position.coords.latitude,
            position.coords.longitude
          );

          mapPosition.setPosition(
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
    mapPosition.setPosition(37.566535, 126.9779692);
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

  bouncy.register();

  return (
    <Styled.Container>
      <MapHeader map={map} />
      <Styled.MapContainer ref={mapRef} />
      {isLoading && (
        <CenterBox bg="black">
          <CenterBox bg="black">
            <l-bouncy size="80" speed="1.75" color="white" />
          </CenterBox>
        </CenterBox>
      )}

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
    </Styled.Container>
  );
};

export default Map;
