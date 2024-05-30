"use client";

import { cn } from "@/lib/utils";
import usePageLoadingStore from "@/store/usePageLoadingStore";
import { useRouter } from "next/navigation";
import GrowBox from "./GrowBox";

interface Props {
  emoji?: string;
  text: string;
  subText?: string;
  center?: boolean;
  url?: string;
  onClickFn?: VoidFunction;
  className?: string;
}

const EmojiHoverButton = ({
  emoji,
  text,
  subText,
  center,
  url,
  onClickFn,
  className,
}: Props) => {
  const router = useRouter();

  const { setLoading } = usePageLoadingStore();

  return (
    <button
      className={cn(
        "block w-full text-left group rounded-sm mb-3 px-1 py-2 hover:bg-black-light-2 mo:text-sm",
        className
      )}
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
