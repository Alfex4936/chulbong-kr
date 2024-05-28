"use client";

import useMarkerImageStore from "@/store/useMarkerImageStore";
import Image from "next/image";
import { useEffect, useRef, useState } from "react";
import { ArrowLeftIcon, ArrowRightIcon } from "../icons/ArrowIcons";
import ExitIcon from "../icons/ExitIcon";
import MinusIcon from "../icons/MinusIcon";
import PlusIcon from "../icons/PlusIcon";

const ImageDetail = () => {
  const { imageView, curImage, nextImage, prevImage, closeImageModal } =
    useMarkerImageStore();

  const [imageSize, setImageSize] = useState(400);

  const outsideRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!imageView) return;

    const handlekeyDown = (e: KeyboardEvent) => {
      if (e.key === "ArrowRight") nextImage();
      if (e.key === "ArrowLeft") prevImage();
    };

    const handleClick = (e: MouseEvent) => {
      if (e.target === outsideRef.current) {
        closeImageModal();
        setImageSize(400);
      }
    };

    window.addEventListener("keydown", handlekeyDown);
    window.addEventListener("click", handleClick);

    return () => {
      window.removeEventListener("keydown", handlekeyDown);
      window.removeEventListener("click", handleClick);
    };
  }, [imageView, outsideRef]);

  const zoomIn = () => {
    if (imageSize < 1000 && window.innerWidth - 20 > imageSize) {
      setImageSize((prev) => prev + 50);
    }
  };

  const zoomOut = () => {
    if (imageSize > 100) setImageSize((prev) => prev - 50);
  };

  if (!imageView || !curImage) return null;

  return (
    <div
      className="absolute top-0 left-0 w-dvw h-dvh bg-black-tp-dark z-[1000] flex justify-center items-center"
      ref={outsideRef}
    >
      <button
        className="absolute top-3 right-3 rounded-full hover:bg-white-tp-light p-1"
        onClick={() => {
          closeImageModal();
          setImageSize(400);
        }}
      >
        <ExitIcon size={25} />
      </button>

      <button
        className="absolute left-24 bg-white-tp-light rounded-full hover:bg-black-tp-light
        mo:left-1/2 mo:-translate-x-[40px] mo:bottom-36"
        onClick={nextImage}
      >
        <ArrowLeftIcon size={30} />
      </button>
      <div className="flex flex-col items-center mx-3">
        <Image
          src={curImage?.photoUrl as string}
          alt="detail"
          width={imageSize}
          height={imageSize}
        />
        <div className="absolute bottom-24 flex botder-red-2 bg-white-tp-light mt-3 px-2 py-1 rounded-lg">
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
      <button
        className="absolute right-24 bg-white-tp-light rounded-full hover:bg-black-tp-light
        mo:right-1/2 mo:translate-x-[40px] mo:bottom-36"
        onClick={prevImage}
      >
        <ArrowRightIcon size={30} />
      </button>
    </div>
  );
};

export default ImageDetail;
