import Image from "next/image";
import { useState } from "react";
import { Skeleton } from "@/components/ui/skeleton";

interface Props {
  src: string;
  w: number;
  h: number;
  alt: string;
}

const ImageWrap = ({ src, h, w, alt }: Props) => {
  const [isLoaded, setIsLoaded] = useState(false);

  return (
    <>
      {!isLoaded && <Skeleton className="w-full h-[140px] mx-auto" />}
      <Image
        src={src}
        width={w}
        height={h}
        alt={alt}
        className={`mx-auto ${isLoaded ? "visible" : "invisible"}`}
        onLoadingComplete={() => setIsLoaded(true)}
      />
    </>
  );
};

export default ImageWrap;
