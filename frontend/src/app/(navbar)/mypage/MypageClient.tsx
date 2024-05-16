"use client";

import EmojiHoverButton from "@/components/atom/EmojiHoverButton";
import ModeToggle from "@/components/common/ModeToggle";
import { Separator } from "@/components/ui/separator";
import useLogout from "@/hooks/mutation/auth/useLogout";
import useMyinfoData from "@/hooks/query/user/useMyinfoData";
import usePageLoadingStore from "@/store/usePageLoadingStore";
import LinkWrap from "./_component/LinkWrap";
import Unauthenticated from "./_component/Unauthenticated";

const MypageClient = () => {
  const { data: myInfo, isError } = useMyinfoData();
  const { setLoading } = usePageLoadingStore();

  const { mutate: logout } = useLogout();

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
        <LinkWrap url="/mypage/user" text="내 정보 관리" />
        <Separator orientation="vertical" className="mx-2 bg-grey-dark-1 h-5" />
        <button className="h-full w-1/2 rounded-md hover:bg-black">설정</button>
      </div>

      <EmojiHoverButton
        emoji="⭐"
        text="저장한 장소"
        subText="북마크 위치"
        url="/mypage/bookmark"
      />
      <EmojiHoverButton
        emoji="🚩"
        text="등록한 장소"
        subText="내가 등록한 위치"
        url="/mypage/mylocate"
      />
      <EmojiHoverButton
        emoji="🪄"
        text="정보 수정 제안 목록"
        subText="내가 수정 제안 한 위치"
        url="/mypage/report"
      />

      <div className="mt-10 mx-auto w-1/2">
        <EmojiHoverButton
          // emoji="🖐️"
          text="로그아웃"
          // subText="다음에 만나요!"
          onClickFn={() => {
            logout();
            setLoading(true);
          }}
          center
        />
      </div>
      <ModeToggle />
      {/* <EmojiHoverButton emoji="🔖📁✏️🚩🗺️⭐❗🖐️✖️🪄🔑" text="저장한 장소" /> */}
    </div>
  );
};

export default MypageClient;
