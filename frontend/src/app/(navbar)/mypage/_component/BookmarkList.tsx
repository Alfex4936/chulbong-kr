import GrowBox from "@/components/atom/GrowBox";
import BookmarkIcon from "@/components/icons/BookmarkIcon";
import { LocationIcon } from "@/components/icons/LocationIcons";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import useDeleteFavorite from "@/hooks/mutation/favorites/useDeleteFavorite";
import useMapStatusStore from "@/store/useMapStatusStore";
import useMapStore from "@/store/useMapStore";
import useMobileMapOpenStore from "@/store/useMobileMapOpenStore";
import { useCallback } from "react";

type Props = {
  title: string;
  subTitle: string;
  lat?: number;
  lng?: number;
  markerId: number;
};

const BookmarkList = ({ title, subTitle, lat, lng, markerId }: Props) => {
  const { open } = useMobileMapOpenStore();
  const { setPosition } = useMapStatusStore();
  const { map, markers } = useMapStore();
  const { mutate: deleteBookmark } = useDeleteFavorite(markerId);

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
  };

  return (
    <li
      className={`flex w-full items-center p-4 rounded-sm mb-2 duration-100 hover:bg-zinc-700 hover:scale-95`}
    >
      <TooltipProvider delayDuration={100}>
        <Tooltip>
          <TooltipTrigger onClick={() => deleteBookmark()}>
            <div className="flex items-center justify-center mr-4 h-8 w-8 rounded-full">
              <BookmarkIcon size={20} isActive />
            </div>
          </TooltipTrigger>
          <TooltipContent>
            <p>북마크 취소</p>
          </TooltipContent>
        </Tooltip>
      </TooltipProvider>

      <div className="w-3/4">
        <div className={`truncate text-left mr-2`}>{title}</div>
        <div className="truncate text-left text-xs text-grey-dark">
          {subTitle}
        </div>
      </div>
      <GrowBox />

      <TooltipProvider delayDuration={100}>
        <Tooltip>
          <TooltipTrigger onClick={filterClickMarker}>
            <div>
              <LocationIcon selected={false} size={18} />
            </div>
          </TooltipTrigger>
          <TooltipContent>
            <p>이동</p>
          </TooltipContent>
        </Tooltip>
      </TooltipProvider>
    </li>
  );
};

export default BookmarkList;
