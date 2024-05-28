"use client";

import useMarkerImageStore from "@/store/useMarkerImageStore";
import Image from "next/image";
import { ArrowLeftIcon, ArrowRightIcon } from "../icons/ArrowIcons";
import ExitIcon from "../icons/ExitIcon";

const ImageDetail = () => {
  const { imageView, curImage, nextImage, prevImage, closeImageModal } =
    useMarkerImageStore();

  if (!imageView || !curImage) return null;

  return (
    <div className="absolute top-0 left-0 w-dvw h-dvh bg-black-tp-dark z-[1000] flex justify-center items-center">
      <button
        className="absolute top-3 right-3 rounded-full hover:bg-white-tp-light p-1"
        onClick={closeImageModal}
      >
        <ExitIcon size={25} />
      </button>

      <button
        className="rounded-full hover:bg-white-tp-light"
        onClick={nextImage}
      >
        <ArrowLeftIcon size={30} />
      </button>
      <div className="mx-3">
        <Image
          src={curImage?.photoUrl as string}
          alt="detail"
          width={400}
          height={300}
        />
      </div>
      <button
        className="rounded-full hover:bg-white-tp-light"
        onClick={prevImage}
      >
        <ArrowRightIcon size={30} />
      </button>
    </div>
  );
};

export default ImageDetail;
