import { MOBILE_WIDTH } from "@/constants";
import useMapStatusStore from "@/store/useMapStatusStore";
import useMapStore from "@/store/useMapStore";
import useMobileMapOpenStore from "@/store/useMobileMapOpenStore";
import usePageLoadingStore from "@/store/usePageLoadingStore";
import { ComponentProps, useCallback } from "react";
import GrowBox from "../atom/GrowBox";
import { LocationIcon } from "../icons/LocationIcons";

interface Props extends ComponentProps<"button"> {
  title: string;
  subTitle?: string;
  lat?: number;
  lng?: number;
  markerId?: number;
  reset: VoidFunction;
}

const MapSearchResult = ({
  title,
  subTitle,
  lat,
  lng,
  markerId,
  reset,
  ...props
}: Props) => {
  const { setLoading } = usePageLoadingStore();

  const { open } = useMobileMapOpenStore();
  const { setPosition } = useMapStatusStore();
  const { map, markers } = useMapStore();

  const { close: mobileMapClose } = useMobileMapOpenStore();

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
    setLoading(true);

    if (window.innerWidth <= MOBILE_WIDTH) {
      mobileMapClose();
    }
  };

  return (
    <button
      className={`flex w-full items-center p-4 rounded-sm mb-2 duration-100 hover:bg-zinc-700 hover:scale-95`}
      onClick={filterClickMarker}
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
