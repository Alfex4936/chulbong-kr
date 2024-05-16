"use client";

import { ComponentProps } from "react";
import GrowBox from "./GrowBox";
import usePageLoadingStore from "@/store/usePageLoadingStore";
import { useRouter } from "next/navigation";

interface Props {
  emoji?: string;
  text: string;
  subText?: string;
  center?: boolean;
  url?: string;
  onClickFn?: VoidFunction;
}

const EmojiHoverButton = ({
  emoji,
  text,
  subText,
  center,
  url,
  onClickFn,
}: Props) => {
  const router = useRouter();

  const { setLoading } = usePageLoadingStore();

  return (
    <button
      className="block w-full text-left group rounded-sm mb-3 px-1 py-2 hover:bg-black-light-2 mo:text-sm"
      onClick={() => {
        if (!url && onClickFn) {
          onClickFn();
          return;
        }
        if (url) {
          router.push(url);
          setLoading(true);
        }
      }}
    >
      <div
        className={`flex ${
          center ? "justify-center" : "justify-start"
        } transition-transform duration-75 transform group-hover:scale-95`}
      >
        {emoji && <span className="mr-2">{emoji}</span>}

        <span>{text}</span>
        {!center && <GrowBox />}
        {subText && (
          <span className="text-sm text-grey-dark-1 mo:text-xs">{subText}</span>
        )}
      </div>
    </button>
  );
};

export default EmojiHoverButton;
