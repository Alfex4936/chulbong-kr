"use client";

import Link from "next/link";
import { ArrowLeftIcon } from "../icons/ArrowIcons";
import usePageLoadingStore from "@/store/usePageLoadingStore";

type Props = {
  url: string;
  text?: string;
};

const PrevHeader = ({ url, text }: Props) => {
  const { setLoading } = usePageLoadingStore();
  return (
    <div className="sticky top-0 left-0 w-full flex items-center h-14 bg-black z-50">
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
};

export default PrevHeader;
