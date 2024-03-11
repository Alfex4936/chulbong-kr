import type { MarkerClusterer } from "@/types/Cluster.types";
import Button from "@mui/material/Button";
import CircularProgress from "@mui/material/CircularProgress";
import { useEffect, useState } from "react";
import activeMarkerImage from "../../assets/images/cb1.webp";
import useInput from "../../hooks/useInput";
import useUploadFormDataStore from "../../store/useUploadFormDataStore";
import type { KakaoMap, KakaoMarker } from "../../types/KakaoMap.types";
import Input from "../Input/Input";
import type { MarkerInfo } from "../Map/Map";
import UploadImage from "../UploadImage/UploadImage";
import * as Styled from "./AddChinupBarForm.style";
import useUploadMarker from "../../hooks/mutation/marker/useUploadMarker";
import { isAxiosError } from "axios";
import logout from "../../api/auth/logout";
import useUserStore from "../../store/useUserStore";

interface Props {
  setState: React.Dispatch<React.SetStateAction<boolean>>;
  setIsMarked: React.Dispatch<React.SetStateAction<boolean>>;
  setMarkerInfoModal: React.Dispatch<React.SetStateAction<boolean>>;
  setCurrentMarkerInfo: React.Dispatch<React.SetStateAction<MarkerInfo | null>>;
  setMarkers: React.Dispatch<React.SetStateAction<KakaoMarker[]>>;
  map: KakaoMap | null;
  marker: KakaoMarker | null;
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
  clusterer,
}: Props) => {
  const formState = useUploadFormDataStore();
  const userState = useUserStore();

  const { mutateAsync: uploadMarker } = useUploadMarker();

  const descriptionValue = useInput("");

  const [error, setError] = useState("");

  const [loading, setLoading] = useState(false);

  useEffect(() => {
    formState.resetData();
  }, []);

  const handleSubmit = async () => {
    const data = {
      description: descriptionValue.value,
      photos: formState.imageForm as File,
      latitude: formState.latitude,
      longitude: formState.longitude,
    };

    setLoading(true);

    try {
      const result = await uploadMarker(data);

      const imageSize = new window.kakao.maps.Size(39, 39);
      const imageOption = { offset: new window.kakao.maps.Point(27, 45) };

      const activeMarkerImg = new window.kakao.maps.MarkerImage(
        activeMarkerImage,
        imageSize,
        imageOption
      );

      const newMarker = new window.kakao.maps.Marker({
        map: map,
        position: new window.kakao.maps.LatLng(
          formState.latitude,
          formState.longitude
        ),
        image: activeMarkerImg,
        title: result.markerId,
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
    } catch (error) {
      if (isAxiosError(error)) {
        if (error.response?.status === 401) {
          await logout();
          userState.resetUser();
          setError("인증이 만료 되었습니다. 다시 로그인 해주세요!");
        } else {
          setError("잠시 후 다시 시도해 주세요!");
        }
      }
      console.log(error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <form>
      <Styled.FormTitle>위치 등록</Styled.FormTitle>

      <UploadImage />

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
    </form>
  );
};

export default AddChinupBarForm;
