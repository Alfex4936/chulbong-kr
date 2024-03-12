import CameraEnhanceIcon from "@mui/icons-material/CameraEnhance";
import PhotoCameraIcon from "@mui/icons-material/PhotoCamera";
import Tooltip from "@mui/material/Tooltip";
import { ChangeEvent, useRef, useState } from "react";
import useUploadFormDataStore from "../../store/useUploadFormDataStore";
import * as Styled from "./UploadImage.tyle";

interface ImageUploadState {
  file: File | null;
  previewURL: string | null;
}

const UploadImage = () => {
  const formState = useUploadFormDataStore();

  const [image, setImage] = useState<ImageUploadState>({
    file: null,
    previewURL: null,
  });
  const [hover, setHover] = useState(false);

  const [errorMessage, setErrorMessage] = useState("");

  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleImageChange = (e: ChangeEvent<HTMLInputElement>) => {
    const suppertedFormats = [
      "image/jpeg",
      "image/png",
      "image/svg+xml",
      "image/webp",
      "image/gif",
    ];
    if (e.target.files) {
      let file = e.target.files[0];
      let reader = new FileReader();

      reader.onloadend = () => {
        setImage({
          file: file,
          previewURL: reader.result as string,
        });
      };

      if (!suppertedFormats.includes(file.type)) {
        setErrorMessage(
          "지원되지 않은 이미지 형식입니다. JPEG, PNG형식의 이미지를 업로드해주세요."
        );
        return;
      }

      if (file.size / (1024 * 1024) > 10) {
        setErrorMessage("이미지는 최대 10MB까지 가능합니다.");
        return;
      }

      setErrorMessage("");

      reader.readAsDataURL(file);

      formState.setImageForm(file);
    }
  };

  const handleBoxClick = () => {
    fileInputRef.current?.click();
  };

  return (
    <Styled.ImageUploadContainer>
      <Tooltip
        title={!image.previewURL ? "사진 등록하기" : "사진 바꾸기"}
        followCursor
      >
        <Styled.ImageBox
          img={image.previewURL}
          onClick={handleBoxClick}
          onMouseEnter={() => {
            setHover(true);
          }}
          onMouseLeave={() => {
            setHover(false);
          }}
          style={{
            border: hover ? "2px dashed #444" : "1px dashed #444",
          }}
        >
          {!image.previewURL ? (
            hover ? (
              <CameraEnhanceIcon
                style={{
                  fontSize: "2rem",
                  color: "#444",
                }}
              />
            ) : (
              <PhotoCameraIcon
                style={{
                  fontSize: "2rem",
                  color: "#444",
                }}
              />
            )
          ) : null}
        </Styled.ImageBox>
      </Tooltip>
      <input type="file" onChange={handleImageChange} ref={fileInputRef} />
      <Styled.ErrorBox>{errorMessage}</Styled.ErrorBox>
    </Styled.ImageUploadContainer>
  );
};

export default UploadImage;
