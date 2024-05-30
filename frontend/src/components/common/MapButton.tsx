"use client";

import useMobileMapOpenStore from "@/store/useMobileMapOpenStore";
import { ComponentProps } from "react";
import { FaMapLocationDot } from "react-icons/fa6";
import MapCloseButton from "./MapCloseButton";

interface Props extends ComponentProps<"button"> {
  selected: boolean;
}

const MapButton = ({ selected, ...props }: Props) => {
  const { toggle, isOpen } = useMobileMapOpenStore();

  return (
    <button
      className="flex justify-center items-center w-12 h-12 shadow-lg bg-black-gradient-1 rounded-[50%]"
      {...props}
      onClick={toggle}
    >
      {isOpen ? <MapCloseButton /> : <FaMapLocationDot className="text-2xl" />}
    </button>
  );
};

export default MapButton;
