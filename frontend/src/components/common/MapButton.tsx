"use client";

import useMobileMapOpenStore from "@/store/useMobileMapOpenStore";
import { ComponentProps } from "react";

interface Props extends ComponentProps<"button"> {
  selected: boolean;
}

const MapButton = ({ selected, ...props }: Props) => {
  const { toggle } = useMobileMapOpenStore();

  return (
    <button
      className="flex justify-center items-center w-12 h-12 shadow-lg bg-black-gradient-1 rounded-[50%]"
      {...props}
      onClick={toggle}
    >
      맵
    </button>
  );
};

export default MapButton;
