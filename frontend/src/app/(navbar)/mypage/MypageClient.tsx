"use client";

import { Separator } from "@/components/ui/separator";
import useLogout from "@/hooks/mutation/auth/useLogout";
import useMyinfoData from "@/hooks/query/user/useMyinfoData";
import Unauthenticated from "./_component/Unauthenticated";
import EmojiHoverButton from "@/components/atom/EmojiHoverButton";
import { useRouter } from "next/navigation";
import ModeToggle from "@/components/common/ModeToggle";
import { useEffect } from "react";

const MypageClient = () => {
  const router = useRouter();

  const { data: myInfo, isError } = useMyinfoData();
  const { mutate: logout } = useLogout();

  useEffect(() => {
    router.prefetch("/mypage/user");
  }, [router]);

  if (!myInfo || isError) return <Unauthenticated />;
  return (
    <div>
      <div className="mb-6">
        <div>
          <span className="text-lg font-bold mo:text-base">
            {myInfo?.username}
          </span>
          님
        </div>
        <div className="text-sm">안녕하세요.</div>
      </div>

      <div className="flex items-center justify-center bg-black-light-2 rounded-md p-1 text-center h-10 mb-6 mo:text-sm">
        <button
          className="h-full w-1/2 rounded-md hover:bg-black"
          onClick={() => router.push("/mypage/user")}
        >
          내 정보 관리
        </button>
        <Separator orientation="vertical" className="mx-2 bg-grey-dark-1 h-5" />
        <button className="h-full w-1/2 rounded-md hover:bg-black">설정</button>
      </div>

      <EmojiHoverButton emoji="⭐" text="저장한 장소" subText="북마크 위치" />
      <EmojiHoverButton
        emoji="🚩"
        text="등록한 장소"
        subText="내가 등록한 위치"
      />

      <div className="mt-10 mx-auto w-1/2">
        <EmojiHoverButton
          // emoji="🖐️"
          text="로그아웃"
          // subText="다음에 만나요!"
          onClick={() => logout()}
          center
        />
      </div>
      <ModeToggle />
      {/* <EmojiHoverButton emoji="🔖📁✏️🚩🗺️⭐❗🖐️✖️🪄🔑" text="저장한 장소" /> */}
    </div>
  );
};

export default MypageClient;
