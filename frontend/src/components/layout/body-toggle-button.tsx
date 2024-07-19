"use client";

import useBodyToggleStore from "@/store/useBodyToggleStore";
import { ArrowLeftIcon, ArrowRightIcon } from "../icons/ArrowIcons";

import IconButton from "../atom/IconButton";

const BodyToggleButton = () => {
  const { isOpen, open } = useBodyToggleStore();

  return (
    <IconButton
      className={`top-1/2 ${
        isOpen ? "left-[500px]" : "left-[80px]"
      } -translate-y-1/2 w-6 h-14 flex justify-center items-center p-0`}
      theme="dark"
      icon={isOpen ? <ArrowLeftIcon /> : <ArrowRightIcon />}
      onClick={open}
    />
  );
};

export default BodyToggleButton;
