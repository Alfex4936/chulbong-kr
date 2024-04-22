"use client";

import { cn } from "@/lib/utils";
import { ArrowLeftIcon, ArrowRightIcon } from "../icons/ArrowIcons";
import useBodyToggleStore from "@/store/useBodyToggleStore";

type Props = {
  children: React.ReactNode;
  toggle?: boolean;
  padding?: boolean;
  containerClass?: string;
  bodyClass?: string;
};

const BlackSideBody = ({
  children,
  toggle,
  containerClass,
  bodyClass,
}: Props) => {
  const { isOpen, open } = useBodyToggleStore();

  return (
    <div
      className={cn(
        `${isOpen ? "web:translate-x-0" : "web:-translate-x-full"} relative ${
          isOpen
            ? "max-w-[410px] mo:min-w-80 min-w-[410px] w-screen"
            : "max-w-[410px] web:w-0 mo:min-w-[320px] mo:w-screen"
        } bg-gradient-to-r from-black to-black-light shadow-lg h-screen
        mo:h-[calc(100%-56px)] mo:m-auto z-10 web:duration-150`,
        containerClass
      )} // web:duration-150
    >
      <div
        className={cn(
          `${
            isOpen ? "web:block" : "web:hidden"
          } px-9 pb-9 h-full overflow-auto scrollbar-thin mo:px-4 mo:pb-4`,
          bodyClass
        )}
      >
        {children}
      </div>
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
