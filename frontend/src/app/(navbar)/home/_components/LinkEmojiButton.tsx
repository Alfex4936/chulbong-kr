"use client";

import GrowBox from "@/components/atom/GrowBox";
import usePageLoadingStore from "@/store/usePageLoadingStore";
import Link from "next/link";

interface Props {
  url: string;
  text: string;
  subText: string;
  emoji: string;
}

const LinkEmojiButton = ({ text, subText, url, emoji }: Props) => {
  const { setLoading } = usePageLoadingStore();
  return (
    <Link href={url} onClick={() => setLoading(true)}>
      <div className="block w-full text-left group rounded-sm mb-3 px-1 py-2 hover:bg-black-light-2 text-sm">
        <div
          className={`flex justify-start transition-transform duration-75 transform group-hover:scale-95`}
        >
          <span className="mr-2">{emoji}</span>

          <span>{text}</span>
          <GrowBox />
          <span className="text-grey-dark-1 text-xs">{subText}</span>
        </div>
      </div>
    </Link>
  );
};

export default LinkEmojiButton;
