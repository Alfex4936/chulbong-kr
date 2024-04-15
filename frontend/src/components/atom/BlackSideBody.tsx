"use client";

import { ArrowLeftIcon, ArrowRightIcon } from "../icons/ArrowIcons";
import useBodyToggleStore from "@/store/useBodyToggleStore";
// TODO: 토글 버튼 색상 변경
// TODO: 토글 애니메이션 다시 적용
// TODO: 그림자 추가

type Props = { children: React.ReactNode; toggle?: boolean };

const BlackSideBody = ({ children, toggle }: Props) => {
  const { isOpen, open } = useBodyToggleStore();

  return (
    <div
      className={`${
        isOpen ? "web:translate-x-0" : "web:-translate-x-full"
      } relative ${
        isOpen
          ? "max-w-[430px] mo:min-w-80 min-w-[430px] w-screen"
          : "max-w-[430px] web:w-0 mo:min-w-[320px] mo:w-screen"
      } bg-gradient-to-r from-black to-black-light mo:h-[100%] mo:m-auto z-10`} // web:duration-150
    >
      <div className={`${isOpen ? "web:block" : "web:hidden"}`}>{children}</div>
      {toggle && (
        <button
          className="absolute -right-9 top-1/2 -translate-y-1/2 
                    bg-gradient-to-r from-black to-black-light py-3 rounded-md shadow-md mo:hidden"
          onClick={() => open()}
        >
          {isOpen ? <ArrowLeftIcon /> : <ArrowRightIcon />}
        </button>
      )}
    </div>
  );
};

export default BlackSideBody;
