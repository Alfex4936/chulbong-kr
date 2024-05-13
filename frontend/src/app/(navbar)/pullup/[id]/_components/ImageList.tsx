import React from "react";
import { type Photo } from "@/types/Marker.types";

type Props = {
  photos?: Photo[];
};

// const images = ["/11.png", "/12.png", "/13.png", "/14.png", "/15.png"];

const ImageList = ({ photos }: Props) => {
  return (
    <div className="flex">
      {photos ? (
        <>
          <div className="w-1/2">
            {photos.map((photo, i) => {
              if (i % 2 === 1) return;
              return (
                <button key={photo.photoId} className="w-full">
                  <img src={photo.photoUrl} className="mx-auto" />
                </button>
              );
            })}
          </div>
          <div className="w-1/2">
            {photos.map((photo, i) => {
              if (i % 2 !== 1) return;
              return (
                <button key={photo.photoId} className="w-full">
                  <img src={photo.photoUrl} className="mx-auto" />
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
