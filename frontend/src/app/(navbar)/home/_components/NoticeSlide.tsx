import GrowBox from "@/components/atom/GrowBox";
import { ArrowRightIcon } from "@/components/icons/ArrowIcons";
import NotificationIcon from "@/components/icons/NotificationIcon";
// TODO: 공지 기능 연동 (후순위)

type Props = {};

const NoticeSlide = (props: Props) => {
  return (
    <button className="flex items-center bg-slate-dark rounded-md p-2 w-full m-auto">
      <span className="mr-2">
        <NotificationIcon selected={false} size={22} />
      </span>
      <span className="truncate text-sm">공지 (레이아웃 용)</span>
      <GrowBox />
      <span>
        <ArrowRightIcon size={22} />
      </span>
    </button>
  );
};

export default NoticeSlide;
