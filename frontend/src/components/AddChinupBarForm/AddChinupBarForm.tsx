import { Button } from "@mui/material";
import { useState } from "react";
import setNewMarker from "../../api/markers/setNewMarker";
import customMarkerImage from "../../assets/images/cb1.png";
import useInput from "../../hooks/useInput";
import useUploadFormDataStore from "../../store/useUploadFormDataStore";
import type { KakaoMap, KakaoMarker } from "../../types/KakaoMap.types";
import Input from "../Input/Input";
import UploadImage from "../UploadImage/UploadImage";
import * as Styled from "./AddChinupBarForm.style";

interface Props {
  setState: React.Dispatch<React.SetStateAction<boolean>>;
  setIsMarked: React.Dispatch<React.SetStateAction<boolean>>;
  map: KakaoMap | null;
  marker: KakaoMarker | null;
}

const AddChinupBarForm = ({ setState, setIsMarked, map, marker }: Props) => {
  const formState = useUploadFormDataStore();

  const descriptionValue = useInput("");

  const [error, setError] = useState("");

  const handleSubmit = () => {
    const data = {
      description: descriptionValue.value,
      photos: formState.imageForm as File,
      latitude: formState.latitude,
      longitude: formState.longitude,
    };

    setNewMarker(data)
      .then((res) => {
        const imageSize = new window.kakao.maps.Size(50, 59);
        const imageOption = { offset: new window.kakao.maps.Point(27, 45) };

        const markerImage2 = new window.kakao.maps.MarkerImage(
          customMarkerImage,
          imageSize,
          imageOption
        );

        new window.kakao.maps.Marker({
          map: map,
          position: new window.kakao.maps.LatLng(
            formState.latitude,
            formState.longitude
          ),
          title: descriptionValue.value,
          image: markerImage2,
        });

        setState(false);
        setIsMarked(false);
        marker?.setMap(null);

        console.log(res);
      })
      .catch((error) => {
        console.log(error);
        setError(error);
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
          backgroundColor: "#333",
          marginTop: "1rem",
          "&:hover": {
            backgroundColor: "#555",
          },
        }}
      >
        등록하기
      </Button>
    </form>
  );
};

export default AddChinupBarForm;
