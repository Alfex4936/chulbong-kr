"use client";

import useMarkerImageStore from "@/store/useMarkerImageStore";
import Image from "next/image";
import { useEffect, useState } from "react";
import { Swiper, SwiperSlide } from "swiper/react";
import LoadingSpinner from "../atom/LoadingSpinner";
import ExitIcon from "../icons/ExitIcon";
import MinusIcon from "../icons/MinusIcon";
import PlusIcon from "../icons/PlusIcon";
import SlideButton from "./SlideButton";

import "swiper/css";
import "swiper/css/navigation";

const ImageDetail = () => {
  const { images, imageView, closeImageModal, curImageIndex } =
    useMarkerImageStore();

  const [imageSize, setImageSize] = useState(400);

  const [isLoaded, setIsLoaded] = useState(false);

  useEffect(() => {
    if (!imageView) return;

    const handlekeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape") closeImageModal();
    };

    window.addEventListener("keydown", handlekeyDown);

    return () => window.removeEventListener("keydown", handlekeyDown);
  }, []);

  const zoomIn = () => {
    if (
      imageSize < 1000 &&
      window.innerWidth - 70 > imageSize &&
      window.innerHeight - 70 > imageSize
    ) {
      setImageSize((prev) => prev + 50);
    }
  };

  const zoomOut = () => {
    if (imageSize > 100) setImageSize((prev) => prev - 50);
  };

  if (!imageView || !images) return null;

  return (
    <div className="absolute top-0 left-0 w-dvw h-dvh bg-black-tp-dark z-[1000] flex justify-center items-center">
      <button
        className="absolute top-3 right-3 rounded-full hover:bg-white-tp-light p-1 z-[1000]"
        onClick={() => {
          closeImageModal();
          setImageSize(400);
        }}
      >
        <ExitIcon size={25} />
      </button>

      <Swiper
        spaceBetween={20}
        slidesPerView={1}
        loop={true}
        initialSlide={curImageIndex}
        className="w-full h-full"
      >
        {images.length > 1 && <SlideButton type="prev" />}
        {images?.map((image) => {
          return (
            <SwiperSlide key={image.photoId}>
              {!isLoaded && (
                <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2">
                  <LoadingSpinner />
                </div>
              )}
              <Image
                src={image.photoUrl}
                alt="detail"
                width={imageSize}
                height={imageSize}
                onLoadingComplete={() => setIsLoaded(true)}
                className={`${
                  isLoaded ? "visible" : "invisible"
                } absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 select-none`}
              />
            </SwiperSlide>
          );
        })}
        {images.length > 1 && <SlideButton type="next" />}
      </Swiper>

      <div className="absolute bottom-24 flex botder-red-2 bg-white-tp-light mt-3 px-2 py-1 rounded-lg z-[1000] mo:hidden">
        <button
          className="hover:bg-black-tp-light rounded-full"
          onClick={zoomOut}
        >
          <MinusIcon />
        </button>
        <div className="mx-1" />
        <button
          className="hover:bg-black-tp-light rounded-full"
          onClick={zoomIn}
        >
          <PlusIcon />
        </button>
      </div>
    </div>
  );
};

export default ImageDetail;
