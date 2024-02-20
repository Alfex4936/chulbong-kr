import { Button } from "@mui/material";
import { useState } from "react";
import setNewMarker from "../../api/markers/setNewMarker";
import useInput from "../../hooks/useInput";
import useUploadFormDataStore from "../../store/useUploadFormDataStore";
import Input from "../Input/Input";
import UploadImage from "../UploadImage/UploadImage";
import * as Styled from "./AddChinupBarForm.style";

const AddChinupBarForm = () => {
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
        console.log(res);
      })
      .catch((error) => {
        console.log(error);
        setError(error);
      });
    console.log(data);
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
