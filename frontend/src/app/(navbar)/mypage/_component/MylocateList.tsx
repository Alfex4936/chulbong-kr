import GrowBox from "@/components/atom/GrowBox";
import LoadingSpinner from "@/components/atom/LoadingSpinner";
import DeleteIcon from "@/components/icons/DeleteIcon";
import { LocationIcon } from "@/components/icons/LocationIcons";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { MAP_LAT_DIF } from "@/constants";
import useDeleteMarker from "@/hooks/mutation/marker/useDeleteMarker";
import useMapStatusStore from "@/store/useMapStatusStore";
import useMapStore from "@/store/useMapStore";
import useMobileMapOpenStore from "@/store/useMobileMapOpenStore";
import usePageLoadingStore from "@/store/usePageLoadingStore";
import { useRouter } from "next/navigation";
import { useCallback, useRef } from "react";

type Props = {
  title: string;
  subTitle: string;
  lat?: number;
  lng?: number;
  markerId: number;
  deleteOption?: boolean;
  isFetching?: boolean;
};

const MylocateList = ({
  title,
  subTitle,
  lat,
  lng,
  markerId,
  deleteOption = true,
  isFetching = false,
}: Props) => {
  const router = useRouter();

  const { setLoading } = usePageLoadingStore();
  const { open } = useMobileMapOpenStore();
  const { setPosition } = useMapStatusStore();
  const { map, markers } = useMapStore();
  const { mutate: deleteMarker, isPending: deletePending } = useDeleteMarker({
    id: markerId,
  });

  const alertRef = useRef<HTMLButtonElement>(null);

  const moveLocation = useCallback(() => {
    const moveLatLon = new window.kakao.maps.LatLng(
      (lat as number) + MAP_LAT_DIF,
      lng
    );
    setPosition((lat as number) + MAP_LAT_DIF, lng as number);
    map?.panTo(moveLatLon);
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
      className={`flex w-full items-center p-4 rounded-sm cursor-pointer mb-2 duration-100 hover:bg-zinc-700 hover:scale-95`}
      onClick={() => {
        setLoading(true);
        const moveLatLon = new window.kakao.maps.LatLng(
          (lat as number) + MAP_LAT_DIF,
          lng
        );

        map?.panTo(moveLatLon);
        router.push(`/pullup/${markerId}`);
      }}
    >
      {deleteOption && (
        <TooltipProvider delayDuration={100}>
          <Tooltip>
            <TooltipTrigger
              onClick={(e) => {
                e.stopPropagation();
                if (deletePending) return;
                if (!alertRef) return;

                alertRef.current?.click();
              }}
              disabled={deletePending || isFetching}
            >
              <div className="flex items-center justify-center mr-4 h-8 w-8 rounded-full">
                {deletePending ? (
                  <LoadingSpinner size="xs" />
                ) : (
                  <DeleteIcon size={20} />
                )}
              </div>
            </TooltipTrigger>
            <TooltipContent>
              <p>삭제</p>
            </TooltipContent>
          </Tooltip>
        </TooltipProvider>
      )}

      <AlertDialog>
        <AlertDialogTrigger asChild>
          <button
            className="hidden"
            ref={alertRef}
            onClick={(e) => {
              e.stopPropagation();
            }}
          >
            마커 삭제
          </button>
        </AlertDialogTrigger>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>정말 삭제하시겠습니까?</AlertDialogTitle>
            <AlertDialogDescription className="text-red">
              저장 된 모든 내용이 사라집니다.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel
              onClick={(e) => {
                e.stopPropagation();
              }}
            >
              취소
            </AlertDialogCancel>
            <AlertDialogAction
              onClick={(e) => {
                e.stopPropagation();
                deleteMarker();
              }}
            >
              삭제
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      <div className="w-3/4">
        <div className={`truncate text-left mr-2`}>{title}</div>
        <div className="truncate text-left text-xs text-grey-dark">
          {subTitle}
        </div>
      </div>
      <GrowBox />

      <TooltipProvider delayDuration={100}>
        <Tooltip>
          <TooltipTrigger
            onClick={(e) => {
              e.stopPropagation();
              filterClickMarker();
            }}
          >
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

export default MylocateList;
