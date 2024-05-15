import { type Photo } from "@/types/Marker.types";
import ImageWrap from "./ImageWrap";

type Props = {
  photos?: Photo[];
};

// const images = ["/11.png", "/12.png", "/13.png", "/14.png", "/15.png"];

const ImageList = ({ photos }: Props) => {
  return (
    <div className="flex">
      {photos ? (
        <>
          <div className="w-1/2 mr-1">
            {photos.map((photo, i) => {
              if (i % 2 === 1) return;
              return (
                <button key={photo.photoId} className="w-full">
                  <ImageWrap src={photo.photoUrl} w={230} h={230} alt="상세" />
                </button>
              );
            })}
          </div>
          <div className="w-1/2 ml-1">
            {photos.map((photo, i) => {
              if (i % 2 !== 1) return;
              return (
                <button key={photo.photoId} className="w-full">
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
