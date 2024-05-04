import useMapStatusStore from "@/store/useMapStatusStore";
import useMapStore from "@/store/useMapStore";
import useMobileMapOpenStore from "@/store/useMobileMapOpenStore";
import { useRouter } from "next/navigation";
import { ComponentProps, useCallback } from "react";
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
  iconClickFn,
  ...props
}: Props) => {
  const router = useRouter();

  const { open } = useMobileMapOpenStore();
  const { setPosition } = useMapStatusStore();
  const { map, markers } = useMapStore();

  const moveLocation = useCallback(() => {
    const moveLatLon = new window.kakao.maps.LatLng(lat, lng);

    setPosition(lat as number, lng as number);
    map?.setCenter(moveLatLon);
    open();
  }, [lat, lng, map]);

  const filterClickMarker = () => {
    if (!markers) return;
    const imageSize = new window.kakao.maps.Size(39, 39);
    const imageOption = { offset: new window.kakao.maps.Point(27, 45) };

    const selectedMarkerImg = new window.kakao.maps.MarkerImage(
      "/selectedMarker.svg",
      imageSize,
      imageOption
    );

    const activeMarkerImg = new window.kakao.maps.MarkerImage(
      "/activeMarker.svg",
      imageSize,
      imageOption
    );

    markers.forEach((marker) => {
      if (Number(marker.getTitle()) === markerId) {
        marker.setImage(selectedMarkerImg);
      } else {
        marker.setImage(activeMarkerImg);
      }
    });

    moveLocation();
    router.push(`pullup/${markerId}`);
  };

  return (
    <button
      className={`flex w-full items-center ${
        styleType === "ranking" ? "p-4" : "p-4"
      } rounded-sm mb-2 duration-100 hover:bg-zinc-700 hover:scale-95`}
      onClick={filterClickMarker}
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
