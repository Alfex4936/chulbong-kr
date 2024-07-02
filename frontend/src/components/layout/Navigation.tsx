"use client";

import NavigationButton from "@/components/atom/NavigationButton";
import ChatBubbleIcon from "@/components/icons/ChatBubbleIcon";
import HomeIcon from "@/components/icons/HomeIcon";
import UserCircleIcon from "@/components/icons/UserCircleIcon";
import Image from "next/image";
import Link from "next/link";
import { usePathname } from "next/navigation";
import MapButton from "../common/MapButton";
import SearchIcon from "../icons/SearchIcon";

const Navigation = () => {
  const pathname = usePathname();

  return (
    <div
      className="w-16 h-screen bg-black-light shadow-xl p-4 z-20 mo:w-screen mo:bottom-0 mo:flex mo:flex-col
                mo:items-center mo:justify-center mo:border-t mo:border-solid mo:p-0 mo:h-[70px]"
    >
      <div
        className="flex flex-col items-center mo:flex-row mo:max-w-[430px] mo:w-full 
                  mo:min-w-80 mo:justify-between mo:ml-auto mo:mr-auto mo:rounded-xl mo:bg-black-light"
      >
        <Link
          href={"/home"}
          className="text-grey text-center w-10 mt-2 mb-6 mo:hidden"
        >
          <Image
            src={"/2.png"}
            width={100}
            height={100}
            alt="logo"
            className="mb-[2px]"
          />
          <p className="text-xs">철봉</p>
        </Link>
        <div className="mo:flex mo:w-1/2 mo:justify-around">
          <NavigationButton
            text="홈"
            url="/home"
            icon={
              <HomeIcon size={25} selected={pathname.startsWith("/home")} />
            }
          />
          <NavigationButton
            text="채팅"
            url="/chat"
            icon={
              <ChatBubbleIcon
                size={25}
                selected={pathname.startsWith("/chat")}
              />
            }
          />
        </div>

        <div className="left-1/2 -translate-y-1/3 web:hidden">
          <MapButton selected />
        </div>

        <div className="mo:flex mo:w-1/2 mo:justify-around">
          <NavigationButton
            text="주변"
            url="/search"
            icon={
              <SearchIcon size={25} selected={pathname.startsWith("/search")} />
            }
          />
          <NavigationButton
            text="내 정보"
            url="/mypage"
            icon={
              <UserCircleIcon
                size={25}
                selected={pathname.startsWith("/mypage")}
              />
            }
          />
        </div>
      </div>
      <div className="w-full h-4 web:hidden" />
    </div>
  );
};

export default Navigation;
