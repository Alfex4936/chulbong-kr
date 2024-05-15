"use client";

import Link from "next/link";
import { ArrowLeftIcon } from "../icons/ArrowIcons";
import usePageLoadingStore from "@/store/usePageLoadingStore";
import { useRouter } from "next/navigation";

type Props = {
  url?: string;
  text?: string;
  back?: boolean;
};

const PrevHeader = ({ url, text, back = false }: Props) => {
  const router = useRouter();

  const { setLoading } = usePageLoadingStore();

  if (back) {
    return (
      <div className="sticky top-0 left-0 w-full flex items-center h-10 bg-gradient-to-r from-black to-black-light z-[300]">
        <button
          className="flex justify-center items-center w-10 h-10 mr-2"
          onClick={() => {
            setLoading(true);
            router.back();
          }}
        >
          <ArrowLeftIcon />
        </button>
        {text && <span>{text}</span>}
      </div>
    );
  }

  if (url) {
    return (
      <div className="sticky top-0 left-0 w-full flex items-center h-10 bg-gradient-to-r from-black to-black-light z-[300]">
        <Link
          href={url}
          className="flex justify-center items-center w-10 h-10 mr-2"
          onClick={() => setLoading(true)}
        >
          <ArrowLeftIcon />
        </Link>
        {text && <span>{text}</span>}
      </div>
    );
  }
};

export default PrevHeader;
