import useMarkerImageStore from "@/store/useMarkerImageStore";
import { type Photo } from "@/types/Marker.types";
import { useEffect } from "react";
import ImageWrap from "./ImageWrap";

type Props = {
  photos?: Photo[];
};

const ImageList = ({ photos }: Props) => {
  const { setImages, setCurImage, openImageModal, setCurImageIndex } =
    useMarkerImageStore();

  useEffect(() => {
    if (!photos) return;

    setImages(photos);
  }, [photos]);

  return (
    <div className="flex">
      {photos ? (
        <>
          <div className="w-1/2 mr-1">
            {photos.map((photo, i) => {
              if (i % 2 === 1) return;
              return (
                <button
                  key={photo.photoId}
                  className="w-full"
                  onClick={() => {
                    setCurImageIndex(i);
                    setCurImage(photo);
                    openImageModal();
                  }}
                >
                  <ImageWrap src={photo.photoUrl} w={230} h={230} alt="상세" />
                </button>
              );
            })}
          </div>
          <div className="w-1/2 ml-1">
            {photos.map((photo, i) => {
              if (i % 2 !== 1) return;
              return (
                <button
                  key={photo.photoId}
                  className="w-full"
                  onClick={() => {
                    setCurImageIndex(i);
                    setCurImage(photo);
                    openImageModal();
                  }}
                >
                  <ImageWrap src={photo.photoUrl} w={230} h={230} alt="상세" />
                </button>
              );
            })}
          </div>
        </>
      ) : (
        <div>등록된 사진이 없습니다.</div>
      )}
    </div>
  );
};

export default ImageList;
