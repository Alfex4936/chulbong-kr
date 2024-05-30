import { MAP_LAT_DIF } from "@/constants";
import useMapStore from "@/store/useMapStore";
import useMiniMapStatusStore from "@/store/useMiniMapStatusStore";
import { ComponentProps, useCallback } from "react";
import GrowBox from "../atom/GrowBox";
import { LocationIcon } from "../icons/LocationIcons";

interface Props extends ComponentProps<"button"> {
  title: string;
  subTitle?: string;
  lat?: number;
  lng?: number;
  mini?: boolean;
  reset: VoidFunction;
  setResultModal: React.Dispatch<React.SetStateAction<boolean>>;
}

const MapSearchResult = ({
  title,
  subTitle,
  lat,
  lng,
  mini = false,
  reset,
  setResultModal,
  ...props
}: Props) => {
  const { map } = useMapStore();
  const { map: minimap } = useMiniMapStatusStore();

  const moveLocation = useCallback(
    (e: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
      e.stopPropagation();
      const moveLatLon = new window.kakao.maps.LatLng(
        (lat as number) + MAP_LAT_DIF,
        lng
      );
      setResultModal(false);

      if (mini) {
        minimap?.panTo(moveLatLon);
        return;
      }
      map?.panTo(moveLatLon);
    },
    [lat, lng, map]
  );

  return (
    <button
      className={`flex w-full items-center p-4 rounded-sm mb-2 duration-100 hover:bg-zinc-700 hover:scale-95`}
      onClick={(e) => moveLocation(e)}
      {...props}
    >
      <div className="w-3/4">
        <div className={`truncate text-left mr-2 text-base`}>{title}</div>
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

export default MapSearchResult;
