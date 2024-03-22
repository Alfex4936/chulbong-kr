import type { MarkerClusterer } from "@/types/Cluster.types";
import AddIcon from "@mui/icons-material/Add";
import RemoveIcon from "@mui/icons-material/Remove";
import Button from "@mui/material/Button";
import CircularProgress from "@mui/material/CircularProgress";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import { isAxiosError } from "axios";
import { useEffect, useState } from "react";
import logout from "../../api/auth/logout";
import activeMarkerImage from "../../assets/images/cb1.webp";
import selectedMarkerImage from "../../assets/images/cb3.webp";
import useSetFacilities from "../../hooks/mutation/marker/useSetFacilities";
import useUploadMarker from "../../hooks/mutation/marker/useUploadMarker";
import useInput from "../../hooks/useInput";
import useCurrentMarkerStore from "../../store/useCurrentMarkerStore";
import useUploadFormDataStore from "../../store/useUploadFormDataStore";
import useUserStore from "../../store/useUserStore";
import type { KakaoMap, KakaoMarker } from "../../types/KakaoMap.types";
import Input from "../Input/Input";
import type { MarkerInfo } from "../Map/Map";
import UploadImage from "../UploadImage/UploadImage";
import * as Styled from "./AddChinupBarForm.style";

interface Props {
  setState: React.Dispatch<React.SetStateAction<boolean>>;
  setIsMarked: React.Dispatch<React.SetStateAction<boolean>>;
  setMarkerInfoModal: React.Dispatch<React.SetStateAction<boolean>>;
  setCurrentMarkerInfo: React.Dispatch<React.SetStateAction<MarkerInfo | null>>;
  setMarkers: React.Dispatch<React.SetStateAction<KakaoMarker[]>>;
  map: KakaoMap | null;
  marker: KakaoMarker | null;
  markers: KakaoMarker[];
  clusterer: MarkerClusterer;
}

const AddChinupBarForm = ({
  setState,
  setIsMarked,
  setMarkerInfoModal,
  setCurrentMarkerInfo,
  setMarkers,
  map,
  marker,
  markers,
  clusterer,
}: Props) => {
  const currentMarkerState = useCurrentMarkerStore();
  const formState = useUploadFormDataStore();
  const userState = useUserStore();

  const { mutateAsync: uploadMarker } = useUploadMarker();

  const [chulbong, setChulbong] = useState(0);
  const [penghang, setPenghang] = useState(0);

  const { mutateAsync: setFacilities } = useSetFacilities();

  const descriptionValue = useInput("");

  const [error, setError] = useState("");

  const [loading, setLoading] = useState(false);

  useEffect(() => {
    formState.resetData();
  }, []);

  const filtering = async () => {
    const imageSize = new window.kakao.maps.Size(39, 39);
    const imageOption = { offset: new window.kakao.maps.Point(27, 45) };

    const activeMarkerImg = new window.kakao.maps.MarkerImage(
      activeMarkerImage,
      imageSize,
      imageOption
    );

    markers.forEach((marker) => {
      marker?.setImage(activeMarkerImg);
    });
  };

  const handleSubmit = async () => {
    const data = {
      description: descriptionValue.value,
      photos: formState.imageForm.map((image) => image.file) as File[],
      latitude: formState.latitude,
      longitude: formState.longitude,
    };

    setLoading(true);

    try {
      const result = await uploadMarker(data);
      await setFacilities({
        markerId: result.markerId,
        facilities: [
          {
            facilityId: 1,
            quantity: chulbong,
          },
          {
            facilityId: 2,
            quantity: penghang,
          },
        ],
      });

      await filtering();

      const imageSize = new window.kakao.maps.Size(39, 39);
      const imageOption = { offset: new window.kakao.maps.Point(27, 45) };

      const selectedMarkerImg = new window.kakao.maps.MarkerImage(
        selectedMarkerImage,
        imageSize,
        imageOption
      );

      const newMarker = new window.kakao.maps.Marker({
        map: map,
        position: new window.kakao.maps.LatLng(
          formState.latitude,
          formState.longitude
        ),
        image: selectedMarkerImg,
        title: result.markerId,
        zIndex: 4,
      });

      window.kakao.maps.event.addListener(newMarker, "click", () => {
        setMarkerInfoModal(true);
        setCurrentMarkerInfo({
          markerId: result.markerId,
        } as MarkerInfo);
      });

      setMarkers((prev) => {
        const copy = [...prev];
        copy.push(newMarker);
        return copy;
      });

      setState(false);
      setIsMarked(false);

      clusterer.addMarker(newMarker);
      marker?.setMap(null);

      currentMarkerState.setMarker(result.markerId);
    } catch (error) {
      if (isAxiosError(error)) {
        if (error.response?.status === 401) {
          await logout();
          userState.resetUser();
          setError("인증이 만료 되었습니다. 다시 로그인 해주세요!");
        } else if (error.response?.status === 409) {
          setError("주변에 이미 철봉이 있습니다!");
        } else if (error.response?.status === 403) {
          setError("대한민국에서만 등록 가능합니다!");
        } else if (error.response?.status === 400) {
          setError("입력을 확인해 주세요!");
        } else {
          setError("잠시 후 다시 시도해 주세요!");
        }
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <form>
      <Styled.FormTitle>위치 등록</Styled.FormTitle>

      <UploadImage />

      <Styled.NumberInputWrap>
        <div>
          <Styled.FlexCenter>철봉</Styled.FlexCenter>
          <Styled.Empty />
          <Styled.FlexCenter>
            <Tooltip title="감소" arrow disableInteractive>
              <IconButton
                aria-label="minus"
                sx={{
                  width: "30px",
                  height: "30px",
                }}
                onClick={() => {
                  if (chulbong === 0) return;
                  setChulbong((prev) => prev - 1);
                }}
              >
                <RemoveIcon />
              </IconButton>
            </Tooltip>
            <Styled.Count>{chulbong}</Styled.Count>
            <Tooltip title="증가" arrow disableInteractive>
              <IconButton
                aria-label="plus"
                sx={{
                  width: "30px",
                  height: "30px",
                }}
                onClick={() => {
                  if (chulbong === 99) return;
                  setChulbong((prev) => prev + 1);
                }}
              >
                <AddIcon />
              </IconButton>
            </Tooltip>
          </Styled.FlexCenter>
        </div>
        <div>
          <Styled.FlexCenter>평행봉</Styled.FlexCenter>
          <Styled.Empty />
          <Styled.FlexCenter>
            <Tooltip title="감소" arrow disableInteractive>
              <IconButton
                aria-label="minus"
                sx={{
                  width: "30px",
                  height: "30px",
                }}
                onClick={() => {
                  if (penghang === 0) return;
                  setPenghang((prev) => prev - 1);
                }}
              >
                <RemoveIcon />
              </IconButton>
            </Tooltip>
            <Styled.Count>{penghang}</Styled.Count>
            <Tooltip title="증가" arrow disableInteractive>
              <IconButton
                aria-label="plus"
                sx={{
                  width: "30px",
                  height: "30px",
                }}
                onClick={() => {
                  if (penghang === 99) return;
                  setPenghang((prev) => prev + 1);
                }}
              >
                <AddIcon />
              </IconButton>
            </Tooltip>
          </Styled.FlexCenter>
        </div>
      </Styled.NumberInputWrap>

      <Styled.InputWrap>
        <Input
          type="text"
          id="description"
          placeholder="설명"
          maxLength={70}
          value={descriptionValue.value}
          onChange={(e) => {
            if (descriptionValue.value.length >= 70) {
              setError("70자 이내로 작성해 주세요!");
            } else {
              setError("");
            }
            descriptionValue.onChange(e);
          }}
        />
        <Styled.ErrorBox>{error}</Styled.ErrorBox>
      </Styled.InputWrap>

      <Button
        onClick={handleSubmit}
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
        disabled={loading}
      >
        {loading ? (
          <CircularProgress size={20} sx={{ color: "#fff" }} />
        ) : (
          "등록하기"
        )}
      </Button>
      {loading && formState.imageForm.length > 0 ? (
        <Styled.ErrorBox>
          이미지를 등록하는 중입니다. 잠시만 기다려 주세요!
        </Styled.ErrorBox>
      ) : null}
    </form>
  );
};

export default AddChinupBarForm;
