"use client";

import { cn } from "@/lib/utils";
import useBodyToggleStore from "@/store/useBodyToggleStore";
import { usePathname } from "next/navigation";
import { useEffect } from "react";
import { ArrowLeftIcon, ArrowRightIcon } from "../icons/ArrowIcons";
// TODO: 모바일 스크롤 멈춤 오류 해결하기
// TODO: 모바일 높이 설정
// TODO: 페이지 전환 로딩바 표시

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
  const pathname = usePathname();
  const { isOpen, open } = useBodyToggleStore();

  useEffect(() => {
    if (!isOpen) {
      open();
    }
  }, [pathname]);

  return (
    <>
      <div
        className={cn(
          `${isOpen ? "web:translate-x-0" : "web:-translate-x-full"} relative ${
            isOpen
              ? "max-w-[410px] mo:min-w-80 min-w-[410px] w-screen"
              : "max-w-[410px] web:w-0 mo:min-w-[320px] mo:w-screen"
          } bg-gradient-to-r from-black to-black-light shadow-lg h-dvh mo:m-auto z-10 web:duration-150`,
          containerClass
        )} // web:duration-150
      >
        <div
          className={cn(
            `${
              isOpen ? "web:block" : "web:hidden"
            } px-9 pb-9 h-full overflow-auto scrollbar-thin mo:px-4 mo:pb-20`,
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

      {isOpen && <div className="absolute top-0 left-0 w-dvw h-dvh bg-black" />}
    </>
  );
};

export default BlackSideBody;
