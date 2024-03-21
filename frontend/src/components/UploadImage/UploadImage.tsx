import CameraEnhanceIcon from "@mui/icons-material/CameraEnhance";
import PhotoCameraIcon from "@mui/icons-material/PhotoCamera";
import Tooltip from "@mui/material/Tooltip";
import { nanoid } from "nanoid";
import { ChangeEvent, useRef, useState } from "react";
import useUploadFormDataStore from "../../store/useUploadFormDataStore";
import resizeFile from "../../utils/resizeFile";
import * as Styled from "./UploadImage.tyle";

export interface ImageUploadState {
  file: File | null;
  previewURL: string | null;
  id: string | null;
}

const UploadImage = () => {
  const formState = useUploadFormDataStore();

  const [images, setImages] = useState<ImageUploadState[]>([]);
  const [hover, setHover] = useState(false);

  const [errorMessage, setErrorMessage] = useState("");

  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleImageChange = async (e: ChangeEvent<HTMLInputElement>) => {
    const suppertedFormats = [
      "image/jpeg",
      "image/png",
      "image/svg+xml",
      "image/webp",
    ];
    if (!e.target.files) return;

    if (!suppertedFormats.includes(e.target.files[0].type)) {
      setErrorMessage(
        "지원되지 않은 이미지 형식입니다. JPEG, PNG, webp형식의 이미지를 업로드해주세요."
      );
      return;
    }

    if (images.length + e.target.files.length > 5) {
      setErrorMessage("최대 5개 까지 등록 가능합니다!");
      return;
    }

    if (e.target.files) {
      let file: File = await resizeFile(e.target.files[0], 0.8);
      let reader = new FileReader();

      reader.onloadend = () => {
        const imageData = {
          file: file,
          previewURL: reader.result as string,
          id: nanoid(),
        };
        setImages((prev) => [...prev, imageData]);
        formState.setImageForm(imageData);
      };

      if (file.size / (1024 * 1024) > 10) {
        setErrorMessage("이미지는 최대 10MB까지 가능합니다.");
        return;
      }

      setErrorMessage("");

      reader.readAsDataURL(file);
    }
  };

  const handleBoxClick = () => {
    fileInputRef.current?.click();
  };

  const deleteImage = (id: string) => {
    const filtered = images.filter((image) => image.id !== id);
    setImages(filtered);
    formState.replaceImages(filtered);
  };

  return (
    <Styled.ImageUploadContainer>
      <div>
        <Tooltip title={"이미지 추가하기"} followCursor>
          <Styled.ImageBox
            img={null}
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
            {hover ? (
              <CameraEnhanceIcon
                style={{
                  color: "#444",
                }}
              />
            ) : (
              <PhotoCameraIcon
                style={{
                  color: "#444",
                }}
              />
            )}
          </Styled.ImageBox>
        </Tooltip>
        <input
          type="file"
          onChange={handleImageChange}
          ref={fileInputRef}
          data-testid="file-input"
        />
        <Styled.ErrorBox data-testid="file-error">
          {errorMessage}
        </Styled.ErrorBox>
      </div>

      <Styled.ImageViewContainer>
        {images.map((image, index) => {
          return (
            <Tooltip title="삭제" arrow disableInteractive key={index}>
              <Styled.ImageView
                img={image.previewURL}
                onClick={() => {
                  deleteImage(image.id as string);
                }}
              />
            </Tooltip>
          );
        })}
      </Styled.ImageViewContainer>
    </Styled.ImageUploadContainer>
  );
};

export default UploadImage;
