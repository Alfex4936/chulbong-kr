import { Button } from "@mui/material";
import Input from "../Input/Input";
import * as Styled from "./AddChinupBarForm.style";
import useInput from "../../hooks/useInput";
import UploadImage from "../UploadImage/UploadImage";
import useUploadFormDataStore from "../../store/useUploadFormDataStore";

const AddChinupBarForm = () => {
  const formState = useUploadFormDataStore();

  const descriptionValue = useInput("");

  const handleSubmit = () => {
    const data = {
      description: descriptionValue.value,
      photoUrl: formState.photoUrl,
      latitude: formState.latitude,
      longitude: formState.longitude,
    };
    console.log(data);
    console.log(formState.imageForm);
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
        {/* <Styled.ErrorBox>{emailError}</Styled.ErrorBox> */}
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
