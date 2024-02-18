import { ChangeEvent, useRef, useState } from "react";
import * as Styled from "./UploadImage.tyle";
import PhotoCameraIcon from "@mui/icons-material/PhotoCamera";
import CameraEnhanceIcon from "@mui/icons-material/CameraEnhance";
import { Tooltip } from "@mui/material";

interface ImageUploadState {
  file: File | null;
  previewURL: string | null;
}

const UploadImage = () => {
  const [image, setImage] = useState<ImageUploadState>({
    file: null,
    previewURL: null,
  });
  const [hover, setHover] = useState(false);

  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleImageChange = (e: ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      let file = e.target.files[0];
      let reader = new FileReader();

      reader.onloadend = () => {
        setImage({
          file: file,
          previewURL: reader.result as string,
        });
      };

      reader.readAsDataURL(file);
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
    </Styled.ImageUploadContainer>
  );
};

export default UploadImage;
