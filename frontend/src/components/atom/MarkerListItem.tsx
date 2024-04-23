import { ComponentProps, useCallback } from "react";
import { LocationIcon } from "../icons/LocationIcons";
import GrowBox from "./GrowBox";
import useMapStore from "@/store/useMapStore";
import useMapStatusStore from "@/store/useMapStatusStore";

interface Props extends ComponentProps<"button"> {
  styleType?: "ranking" | "normal";
  title: string;
  subTitle?: string;
  ranking?: number;
  lat?: number;
  lng?: number;
}

const MarkerListItem = ({
  styleType = "normal",
  title,
  subTitle,
  ranking,
  lat,
  lng,
  ...props
}: Props) => {
  const { setPosition } = useMapStatusStore();
  const { map } = useMapStore();

  const moveLocation = useCallback(() => {
    const moveLatLon = new window.kakao.maps.LatLng(lat, lng);

    setPosition(lat as number, lng as number);
    map?.setCenter(moveLatLon);
  }, [lat, lng, map]);

  return (
    <button
      className={`flex w-full items-center ${
        styleType === "ranking" ? "p-4" : "p-1"
      } rounded-sm mb-2 duration-100 hover:bg-zinc-700 hover:scale-95`}
      onClick={moveLocation}
      {...props}
    >
      {styleType === "ranking" && (
        <div className="mr-4">
          {ranking}
          <span className="text-xs text-grey-dark">등</span>
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
