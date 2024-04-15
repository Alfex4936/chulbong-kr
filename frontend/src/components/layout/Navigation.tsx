"use client";

import IconButton from "@/components/atom/IconButton";
import ChatBubbleIcon from "@/components/icons/ChatBubbleIcon";
import HomeIcon from "@/components/icons/HomeIcon";
// import { LocationIcon } from "@/components/icons/LocationIcons";
import NotificationIcon from "@/components/icons/NotificationIcon";
// import SettingIcon from "@/components/icons/SettingIcon";
import UserCircleIcon from "@/components/icons/UserCircleIcon";
import { usePathname, useRouter } from "next/navigation";
import MapButton from "../common/MapButton";

const Navigation = () => {
  const pathname = usePathname();
  const router = useRouter();

  const pushRouter = (url: string) => {
    router.push(url);
  };

  return (
    <div
      className="w-16 h-screen bg-black-light shadow-mdxl p-4 z-20
                    mo:w-screen mo:h-14 mo:flex mo:items-center mo:pl-0 mo:pr-0 mo:border-t mo:border-solid mo:border-black"
    >
      <div
        className="flex flex-col items-center 
                      mo:flex-row mo:max-w-[430px] mo:w-full mo:min-w-80 mo:justify-between mo:ml-auto mo:mr-auto"
      >
        <div className="text-grey text-center w-10 mt-2 mb-6 mo:hidden">
          철봉
        </div>
        <div className="mo:flex mo:w-1/2 mo:justify-around">
          <IconButton
            text="홈"
            selected={pathname === "/home"}
            icon={<HomeIcon size={25} selected={pathname === "/home"} />}
            onClick={() => pushRouter("/home")}
          />
          <IconButton
            text="채팅"
            selected={pathname === "/chat"}
            icon={<ChatBubbleIcon size={25} selected={pathname === "/chat"} />}
            onClick={() => pushRouter("/chat")}
          />
        </div>

        <div className="left-1/2 -translate-y-1/2 web:hidden">
          <MapButton selected />
        </div>

        <div className="mo:flex mo:w-1/2 mo:justify-around">
          <IconButton
            text="공지"
            selected={pathname === "/notice"}
            icon={
              <NotificationIcon size={25} selected={pathname === "/notice"} />
            }
            onClick={() => pushRouter("/notice")}
          />
          <IconButton
            text="내 정보"
            selected={pathname === "/mypage"}
            icon={
              <UserCircleIcon size={25} selected={pathname === "/mypage"} />
            }
            onClick={() => pushRouter("/mypage")}
          />
        </div>

        {/* <IconButton
        selected={pathname === "/setting"}
        icon={<SettingIcon selected={pathname === "/setting"} />}
        onClick={() => pushRouter("/setting")}
        /> */}
      </div>
    </div>
  );
};

export default Navigation;
