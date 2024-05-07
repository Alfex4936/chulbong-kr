"use client";

import { Button } from "@/components/ui/button";
import useUploadFormDataStore, {
  type ImageUploadState,
} from "@/store/useUploadFormDataStore";
import resizeFile from "@/utils/resizeFile";
import { ChangeEvent, useRef, useState } from "react";
import { FaCameraRetro } from "react-icons/fa";
import { FaCamera } from "react-icons/fa6";
import { v4 } from "uuid";

const UploadImage = () => {
  const { setImageForm, replaceImages } = useUploadFormDataStore();

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

    if (!suppertedFormats.includes(e.target.files[0]?.type)) {
      setErrorMessage(
        "지원되지 않은 이미지 형식입니다. JPEG, PNG, webp형식의 이미지를 업로드해주세요."
      );
      return;
    }

    if (images.length + e.target.files.length > 5) {
      setErrorMessage("최대 5개 까지 등록 가능합니다!");
      return;
    }

    let file: File = await resizeFile(e.target.files[0], 0.8);
    let reader = new FileReader();

    reader.onloadend = () => {
      const imageData = {
        file: file,
        previewURL: reader.result as string,
        id: v4(),
      };
      setImages((prev) => [...prev, imageData]);
      setImageForm(imageData);
    };

    if (file.size / (1024 * 1024) > 10) {
      setErrorMessage("이미지는 최대 10MB까지 가능합니다.");
      return;
    }

    setErrorMessage("");

    reader.readAsDataURL(file);
  };

  const handleBoxClick = () => {
    fileInputRef.current?.click();
  };

  const deleteImage = (id: string) => {
    const filtered = images.filter((image) => image.id !== id);
    setImages(filtered);
    replaceImages(filtered);
  };

  // TODO: 사진기 아이콘
  return (
    <div className="flex flex-col justify-center mb-4">
      <div>
        <div
          className="relative rounded-sm flex items-center justify-center h-14 cursor-pointer m-auto"
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
          {hover ? <FaCameraRetro /> : <FaCamera />}
        </div>
        <input
          type="file"
          onChange={handleImageChange}
          ref={fileInputRef}
          data-testid="file-input"
          className="hidden"
        />
      </div>

      <div className="mt-2 flex justify-center">
        {images.map((image, index) => {
          return (
            <div
              key={index}
              className="flex flex-col items-center justify-center"
            >
              <img
                className="w-14 h-14 m-1 object-cover bg-center bg-no-repeat bg-cover rounded-sm"
                src={image.previewURL as string}
              />
              <Button
                size="sm"
                className="border-grey border bg-transparent dark:text-grey hover:bg-white-tp-light hover:border-transparent"
                onClick={() => {
                  deleteImage(image.id as string);
                }}
              >
                삭제
              </Button>
            </div>
          );
        })}
      </div>

      <div
        data-testid="file-error"
        className="mt-1 text-center text-sm text-red"
      >
        {errorMessage}
      </div>
    </div>
  );
};

export default UploadImage;
