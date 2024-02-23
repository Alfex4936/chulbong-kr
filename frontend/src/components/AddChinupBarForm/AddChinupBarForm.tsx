import { Button, CircularProgress } from "@mui/material";
import { useState } from "react";
import setNewMarker from "../../api/markers/setNewMarker";
import customMarkerImage from "../../assets/images/cb1.png";
import useInput from "../../hooks/useInput";
import useUploadFormDataStore from "../../store/useUploadFormDataStore";
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
  markers: KakaoMarker[];
  marker: KakaoMarker | null;
}

const AddChinupBarForm = ({
  setState,
  setIsMarked,
  setMarkerInfoModal,
  setCurrentMarkerInfo,
  setMarkers,
  map,
  markers,
  marker,
}: Props) => {
  const formState = useUploadFormDataStore();

  const descriptionValue = useInput("");

  const [error, setError] = useState("");

  const [loading, setLoading] = useState(false);

  const handleSubmit = () => {
    const data = {
      description: descriptionValue.value,
      photos: formState.imageForm as File,
      latitude: formState.latitude,
      longitude: formState.longitude,
    };
    setLoading(true);

    setNewMarker(data)
      .then((res) => {
        const imageSize = new window.kakao.maps.Size(50, 59);
        const imageOption = { offset: new window.kakao.maps.Point(27, 45) };

        const markerImage2 = new window.kakao.maps.MarkerImage(
          customMarkerImage,
          imageSize,
          imageOption
        );

        const newMarker = new window.kakao.maps.Marker({
          map: map,
          position: new window.kakao.maps.LatLng(
            formState.latitude,
            formState.longitude
          ),
          title: descriptionValue.value,
          image: markerImage2,
        });

        window.kakao.maps.event.addListener(newMarker, "click", () => {
          setMarkerInfoModal(true);
          setCurrentMarkerInfo({
            ...res?.data,
            index: markers.length,
            userId: res?.data.userID,
            photos: res?.data.photoUrls,
          });
        });

        setMarkers((prev) => {
          const copy = [...prev];
          copy.push(newMarker);
          return copy;
        });

        setState(false);
        setIsMarked(false);
        marker?.setMap(null);
      })
      .catch((error) => {
        console.log(error);
        setError(error);
      })
      .finally(() => {
        setLoading(false);
      });
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
          value={descriptionValue.value}
          onChange={(e) => {
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
