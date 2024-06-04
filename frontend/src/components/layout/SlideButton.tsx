import { useEffect } from "react";
import { useSwiper } from "swiper/react";
import { ArrowLeftIcon, ArrowRightIcon } from "../icons/ArrowIcons";

interface Props {
  type: "next" | "prev";
}

const SlideButton = ({ type }: Props) => {
  const swiper = useSwiper();

  useEffect(() => {
    const handlekeyDown = (e: KeyboardEvent) => {
      if (e.key === "ArrowRight") swiper.slideNext();
      if (e.key === "ArrowLeft") swiper.slidePrev();
    };

    window.addEventListener("keydown", handlekeyDown);

    return () => window.removeEventListener("keydown", handlekeyDown);
  }, []);

  if (type === "next") {
    return (
      <button
        className="absolute top-1/2 -translate-y-1/2 right-4 z-[1000] bg-white-tp-light rounded-full hover:bg-black-tp-light
                    mo:left-1/2 mo:-translate-x-[40px] mo:bottom-36 mo:hidden"
        onClick={() => swiper.slideNext()}
      >
        <ArrowRightIcon size={40} />
      </button>
    );
  }

  if (type === "prev") {
    return (
      <button
        className="absolute top-1/2 -translate-y-1/2 left-4 z-[1000] bg-white-tp-light rounded-full hover:bg-black-tp-light
        mo:right-1/2 mo:translate-x-[40px] mo:bottom-36 mo:hidden"
        onClick={() => swiper.slidePrev()}
      >
        <ArrowLeftIcon size={40} />
      </button>
    );
  }
};

export default SlideButton;
