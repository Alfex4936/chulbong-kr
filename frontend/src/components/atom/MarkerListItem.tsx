import useMapStore from "@/store/useMapStore";
import usePageLoadingStore from "@/store/usePageLoadingStore";
import { useRouter } from "next/navigation";
import { ComponentProps } from "react";
import { LocationIcon } from "../icons/LocationIcons";
import GrowBox from "./GrowBox";

interface Props extends ComponentProps<"button"> {
  styleType?: "ranking" | "normal";
  title: string;
  subTitle?: string;
  ranking?: number;
  lat?: number;
  lng?: number;
  markerId?: number;
  icon?: string | React.ReactNode;
  mini?: boolean;
  searchToggle?: boolean;
  reset?: VoidFunction;
  iconClickFn?: VoidFunction;
}

const MarkerListItem = ({
  styleType = "normal",
  title,
  subTitle,
  ranking,
  lat,
  lng,
  markerId,
  icon,
  mini = false,
  searchToggle,
  reset,
  iconClickFn,
  ...props
}: Props) => {
  const router = useRouter();

  const { setLoading } = usePageLoadingStore();
  const { map } = useMapStore();

  return (
    <button
      className={`flex w-full items-center ${
        styleType === "ranking" ? "p-4" : "p-4"
      } rounded-sm mb-2 duration-100 hover:bg-zinc-700 hover:scale-95`}
      onClick={() => {
        setLoading(true);
        const moveLatLon = new window.kakao.maps.LatLng(
          (lat as number) + 0.003,
          lng
        );

        map?.panTo(moveLatLon);
        router.push(`pullup/${markerId}`);
      }}
      {...props}
    >
      {styleType === "ranking" && (
        <div className="mr-4">
          {ranking}
          <span className="text-xs text-grey-dark">등</span>
        </div>
      )}

      {icon && (
        <div
          className="flex items-center justify-center mr-4 h-8 w-8 rounded-full"
          onClick={(e) => {
            if (!iconClickFn) return;
            e.stopPropagation();
            iconClickFn();
          }}
        >
          {icon}
        </div>
      )}

      <div className="w-3/4">
        <div
          className={`truncate text-left mr-2 ${
            styleType === "ranking" ? "text-sm" : "text-base"
          }`}
        >
          {title}
        </div>
        <div className="truncate text-left text-xs text-grey-dark">
          {subTitle}
        </div>
      </div>
      <GrowBox />
      <div>
        <LocationIcon selected={false} size={18} />
      </div>
    </button>
  );
};

export default MarkerListItem;
