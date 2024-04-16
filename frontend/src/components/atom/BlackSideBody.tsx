"use client";

import { ArrowLeftIcon, ArrowRightIcon } from "../icons/ArrowIcons";
import useBodyToggleStore from "@/store/useBodyToggleStore";

type Props = { children: React.ReactNode; toggle?: boolean };

const BlackSideBody = ({ children, toggle }: Props) => {
  const { isOpen, open } = useBodyToggleStore();

  return (
    <div
      className={`${
        isOpen ? "web:translate-x-0" : "web:-translate-x-full"
      } relative ${
        isOpen
          ? "max-w-[410px] mo:min-w-80 min-w-[410px] w-screen"
          : "max-w-[410px] web:w-0 mo:min-w-[320px] mo:w-screen"
      } bg-gradient-to-r from-black to-black-light shadow-lg px-9 mo:h-[100%] mo:m-auto z-10 web:duration-150`} // web:duration-150
    >
      <div className={`${isOpen ? "web:block" : "web:hidden"}`}>{children}</div>
      {toggle && (
        <button
          className="absolute -right-9 top-1/2 -translate-y-1/2 
                  bg-black-light py-3 rounded-md shadow-md mo:hidden"
          onClick={() => open()}
        >
          {isOpen ? <ArrowLeftIcon /> : <ArrowRightIcon />}
        </button>
      )}
    </div>
  );
};

export default BlackSideBody;
