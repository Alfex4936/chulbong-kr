"use client";

import { cn } from "@/lib/utils";
import useBodyToggleStore from "@/store/useBodyToggleStore";
import usePageLoadingStore from "@/store/usePageLoadingStore";
import useScrollButtonStore from "@/store/useScrollButtonStore";
import { usePathname } from "next/navigation";
import { useEffect, useRef } from "react";
import { FaArrowUp } from "react-icons/fa";
import { ArrowLeftIcon, ArrowRightIcon } from "../icons/ArrowIcons";
import LoadingSpinner from "./LoadingSpinner";

type Props = {
  children: React.ReactNode;
  className?: string;
};

const BlackSideBody = ({ children, className }: Props) => {
  const pathname = usePathname();

  const { isActive, setIsActive } = useScrollButtonStore();
  const { isOpen, open } = useBodyToggleStore();
  const { isLoading, setLoading, setVisible } = usePageLoadingStore();

  const bodyRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    setLoading(false);
    setIsActive(false);
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

  const scrollToTop = () => {
    if (!bodyRef.current) return;
    bodyRef.current.scrollTo({
      top: 0,
      behavior: "smooth",
    });
  };

  return (
    <div
      className={cn(
        `${
          isOpen
            ? "web:translate-x-0 min-w-[410px] w-screen mo:min-w-80"
            : "web:-translate-x-[200%] web:w-0 mo:min-w-[320px] mo:w-screen"
        } web:max-w-[410px] mo:w-full bg-gradient-to-r from-black to-black-light 
          shadow-lg z-10 web:duration-150 h-full mo:h-[calc(100%-70px)] overflow-auto scrollbar-thin`,
        className
      )}
      ref={bodyRef}
    >
      {isLoading ? (
        <div className="h-full flex justify-center items-center">
          <LoadingSpinner />
        </div>
      ) : (
        <>
          {children}

          {isActive && (
            <button
              className="sticky left-[400px] shadow-lg bg-black-gradient-1 flex items-center justify-center 
              w-8 h-8 border-[0.5px] border-gray-900 rounded-full animate-bottom-top"
              onClick={scrollToTop}
            >
              <FaArrowUp />
            </button>
          )}
        </>
      )}
    </div>
  );
};

export default BlackSideBody;
