import CameraEnhanceIcon from "@mui/icons-material/CameraEnhance";
import PhotoCameraIcon from "@mui/icons-material/PhotoCamera";
import Tooltip from "@mui/material/Tooltip";
import { ChangeEvent, useRef, useState } from "react";
import useUploadFormDataStore from "../../store/useUploadFormDataStore";
import * as Styled from "./UploadImage.tyle";
import { nanoid } from "nanoid";

interface ImageUploadState {
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

  const resizeFile = async (file: File, scale: number): Promise<File> => {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = (e) => {
        const img = new Image();
        img.onload = () => {
          const width = img.width * scale;
          const height = img.height * scale;

          const canvas = document.createElement("canvas");
          canvas.width = width;
          canvas.height = height;

          const ctx = canvas.getContext("2d");
          ctx?.drawImage(img, 0, 0, width, height);
          canvas.toBlob((blob) => {
            if (blob) {
              const resizedFile = new File([blob], file.name, {
                type: file.type,
                lastModified: Date.now(),
              });
              resolve(resizedFile);
            } else {
              reject(new Error("Canvas toBlob failed"));
            }
          }, file.type);
        };
        img.src = e.target?.result as string;
      };
      reader.readAsDataURL(file);
    });
  };

  const handleImageChange = async (e: ChangeEvent<HTMLInputElement>) => {
    if (images.length > 4) {
      setErrorMessage("최대 5개 까지 등록 가능합니다!");
      return;
    }
    const suppertedFormats = [
      "image/jpeg",
      "image/png",
      "image/svg+xml",
      "image/webp",
    ];
    if (e.target.files) {
      let file: File = await resizeFile(e.target.files[0], 0.8);
      let reader = new FileReader();

      reader.onloadend = () => {
        setImages((prev) => [
          ...prev,
          {
            file: file,
            previewURL: reader.result as string,
            id: nanoid(),
          },
        ]);
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

  const deleteImage = (id: string) => {
    const filtered = images.filter((image) => image.id !== id);
    setImages(filtered);
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
        <input type="file" onChange={handleImageChange} ref={fileInputRef} />
        <Styled.ErrorBox>{errorMessage}</Styled.ErrorBox>
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
      {/* {images.map((image, index) => {
        return (
          <Fragment key={index}>
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
            <input
              type="file"
              onChange={handleImageChange}
              ref={fileInputRef}
            />
            <Styled.ErrorBox>{errorMessage}</Styled.ErrorBox>
          </Fragment>
        );
      })} */}
    </Styled.ImageUploadContainer>
  );
};

export default UploadImage;
