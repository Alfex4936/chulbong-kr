"use client";

import useBodyToggleStore from "@/store/useBodyToggleStore";
import useMobileMapOpenStore from "@/store/useMobileMapOpenStore";
import usePageLoadingStore from "@/store/usePageLoadingStore";
import Link from "next/link";
import { usePathname } from "next/navigation";

interface Props {
  icon: React.ReactNode;
  text?: string;
  url: string;
}

const IconButton = ({ icon, text, url }: Props) => {
  const pathname = usePathname();

  const { isOpen: isBodyOpen, open: bodyOpen } = useBodyToggleStore();
  const { setLoading } = usePageLoadingStore();
  const { close } = useMobileMapOpenStore();

  return (
    <Link
      href={url}
      className="flex justify-center items-center web:mb-3 mo:w-full"
      onClick={() => {
        if (!isBodyOpen) bodyOpen();
        if (pathname === url) {
          close();
          return;
        }
        setLoading(true);
        close();
      }}
    >
      <div className="flex flex-col justify-center items-center w-12 h-16 cursor-pointer rounded-sm mo:w-[44px] mo:h-[50px]">
        <div className="ml-[1px]">{icon}</div>
        {text && (
          <div
            className={` text-sm mp:font-xs mo:text-[10px] web:mt-2 text-grey`}
          >
            {text}
          </div>
        )}
      </div>
    </Link>
  );
};

export default IconButton;
