"use client";

import { ArrowLeftIcon } from "../icons/ArrowIcons";
import { useRouter } from "next/navigation";

type Props = {
  url: string;
  text?: string;
};

const PrevHeader = ({ url, text }: Props) => {
  const router = useRouter();
  return (
    <div className="flex items-center h-14">
      <button
        className="flex justify-center items-center w-10 h-10 mr-2"
        onClick={() => router.push(url)}
      >
        <ArrowLeftIcon />
      </button>
      {text && <span>{text}</span>}
    </div>
  );
};

export default PrevHeader;
