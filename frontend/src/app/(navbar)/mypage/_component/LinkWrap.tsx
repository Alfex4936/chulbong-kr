"use client";

import usePageLoadingStore from "@/store/usePageLoadingStore";
import Link from "next/link";

type Props = {
  url: string;
  text: string;
};

const LinkWrap = ({ url, text }: Props) => {
  const { setLoading } = usePageLoadingStore();

  return (
    <Link
      href={url}
      className="flex items-center justify-center h-full w-1/2 rounded-md hover:bg-black"
      onClick={() => setLoading(true)}
    >
      {text}
    </Link>
  );
};

export default LinkWrap;
