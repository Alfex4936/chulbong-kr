"use client";

import { cn } from "@/lib/utils";
import useBodyToggleStore from "@/store/useBodyToggleStore";
import usePageLoadingStore from "@/store/usePageLoadingStore";
import { usePathname } from "next/navigation";
import { useEffect } from "react";
import { ArrowLeftIcon, ArrowRightIcon } from "../icons/ArrowIcons";
import LoadingSpinner from "./LoadingSpinner";
// import PageLoadingBar from "../layout/PageLoadingBar";

type Props = {
  children: React.ReactNode;
  toggle?: boolean;
  padding?: boolean;
  bodyClass?: string;
  className?: string;
};

const BlackSideBody = ({ children, toggle, bodyClass, className }: Props) => {
  const pathname = usePathname();

  const { isOpen, open } = useBodyToggleStore();
  const { isLoading, setLoading, setVisible } = usePageLoadingStore();

  useEffect(() => {
    setLoading(false);
  }, []);

  useEffect(() => {
    const time = setTimeout(() => {
      setVisible(false);
    }, 300);

    return () => clearTimeout(time);
  }, [isLoading]);

  useEffect(() => {
    if (!isOpen) {
      open();
    }
  }, [pathname]);

  return (
    <div
      className={cn(
        `${isOpen ? "web:translate-x-0" : "web:-translate-x-[150%]"} relative ${
          isOpen
            ? "mo:min-w-80 min-w-[410px] w-screen"
            : "web:w-0 mo:min-w-[320px] mo:w-screen"
        } web:max-w-[410px] mo:w-full bg-gradient-to-r from-black to-black-light 
          shadow-lg mo:m-auto z-10 web:duration-150 h-full`,
        className
      )}
    >
      {isLoading ? (
        <div className="h-full flex justify-center items-center">
          <LoadingSpinner />
        </div>
      ) : (
        <>
          <div
            className={cn(
              `h-full overflow-auto scrollbar-thin pb-4 mo:pb-20 mo:px-4 ${
                isOpen ? "px-4" : "px-0"
              }`,
              bodyClass
            )}
          >
            {/* <PageLoadingBar /> */}
            {children}
          </div>

          {toggle && (
            <button
              className="absolute -right-9 top-1/2 -translate-y-1/2 z-50
            bg-black-light py-3 rounded-md shadow-md mo:hidden"
              onClick={() => open()}
            >
              {isOpen ? <ArrowLeftIcon /> : <ArrowRightIcon />}
            </button>
          )}
        </>
      )}
    </div>
  );
};

export default BlackSideBody;
