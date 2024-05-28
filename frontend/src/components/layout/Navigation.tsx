"use client";

import IconButton from "@/components/atom/IconButton";
import ChatBubbleIcon from "@/components/icons/ChatBubbleIcon";
import HomeIcon from "@/components/icons/HomeIcon";
import UserCircleIcon from "@/components/icons/UserCircleIcon";
import Image from "next/image";
import Link from "next/link";
import { usePathname } from "next/navigation";
import MapButton from "../common/MapButton";
import SearchIcon from "../icons/SearchIcon";
// import { LocationIcon } from "@/components/icons/LocationIcons";
// import NotificationIcon from "@/components/icons/NotificationIcon";
// import SettingIcon from "@/components/icons/SettingIcon";

const Navigation = () => {
  const pathname = usePathname();

  return (
    <div
      className="w-16 h-screen bg-black-light shadow-xl p-4 z-20 mo:w-screen mo:h-14 mo:fixed mo:bottom-0 mo:flex 
                mo:items-center mo:pl-0 mo:pr-0 mo:border-t mo:border-solid"
    >
      <div
        className="flex flex-col items-center mo:flex-row mo:max-w-[430px] pb-9 pt-2 mo:w-full 
                  mo:min-w-80 mo:justify-between mo:ml-auto mo:mr-auto mo:rounded-xl mo:bg-black-light"
      >
        <Link
          href={"/"}
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
          <IconButton
            text="홈"
            url="/home"
            icon={
              <HomeIcon size={25} selected={pathname.startsWith("/home")} />
            }
          />
          <IconButton
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
          {/* <IconButton
            text="공지"
            url="/notice"
            icon={
              <NotificationIcon
                size={25}
                selected={pathname.startsWith("/notice")}
              />
            }
          /> */}
          <IconButton
            text="주변"
            url="/search"
            icon={
              <SearchIcon size={25} selected={pathname.startsWith("/search")} />
            }
          />
          <IconButton
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
    </div>
  );
};

export default Navigation;
