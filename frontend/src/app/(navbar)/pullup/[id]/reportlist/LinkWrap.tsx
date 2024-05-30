"use client";

import usePageLoadingStore from "@/store/usePageLoadingStore";
import Link from "next/link";

interface Props {
  id: string;
}

const LinkWrap = ({ id }: Props) => {
  const { setLoading } = usePageLoadingStore();

  return (
    <Link
      href={`/pullup/${id}/report`}
      className="flex w-[90%] m-auto text-left group rounded-sm mb-3 px-1 py-2 bg-black-light-2 mo:text-sm"
      onClick={() => setLoading(true)}
    >
      <div>✏️</div>
      <div>제안 요청 하러 가기</div>
    </Link>
  );
};

export default LinkWrap;
