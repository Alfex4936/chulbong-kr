import { MOBILE_WIDTH } from "@/constants";
import useMapStatusStore from "@/store/useMapStatusStore";
import useMapStore from "@/store/useMapStore";
import useMiniMapStatusStore from "@/store/useMiniMapStatusStore";
import useMobileMapOpenStore from "@/store/useMobileMapOpenStore";
import usePageLoadingStore from "@/store/usePageLoadingStore";
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

  const { open } = useMobileMapOpenStore();
  const { setPosition } = useMapStatusStore();
  const { map, markers } = useMapStore();

  const { close: mobileMapClose } = useMobileMapOpenStore();

  const { map: miniMap } = useMiniMapStatusStore();

  const moveLocation = useCallback(() => {
    if (mini) {
      const moveLatLon = new window.kakao.maps.LatLng(lat, lng);
      miniMap?.setCenter(moveLatLon);
    } else {
      const moveLatLon = new window.kakao.maps.LatLng(lat, lng);

      setPosition(lat as number, lng as number);
      map?.setCenter(moveLatLon);
      open();
    }
  }, [lat, lng, map, mini]);

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

    if (window.innerWidth <= MOBILE_WIDTH) {
      mobileMapClose();
    }
  };

  return (
    <button
      className={`flex w-full items-center ${
        styleType === "ranking" ? "p-4" : "p-4"
      } rounded-sm mb-2 duration-100 hover:bg-zinc-700 hover:scale-95`}
      onClick={(e) => {
        e.stopPropagation();
        if (searchToggle) {
          if (!reset) return;
          reset();
        }
        if (mini) {
          moveLocation();
        } else {
          setLoading(true);
          filterClickMarker();
        }
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
