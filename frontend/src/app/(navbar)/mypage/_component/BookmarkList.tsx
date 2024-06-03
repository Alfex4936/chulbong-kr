import GrowBox from "@/components/atom/GrowBox";
import LoadingSpinner from "@/components/atom/LoadingSpinner";
import BookmarkIcon from "@/components/icons/BookmarkIcon";
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
import useMapControl from "@/hooks/common/useMapControl";
import useDeleteFavorite from "@/hooks/mutation/favorites/useDeleteFavorite";
import usePageLoadingStore from "@/store/usePageLoadingStore";
import { useRouter } from "next/navigation";
import { useRef } from "react";

type Props = {
  title: string;
  subTitle: string;
  lat?: number;
  lng?: number;
  markerId: number;
  isFetching?: boolean;
};

const BookmarkList = ({
  title,
  subTitle,
  lat,
  lng,
  markerId,
  isFetching = false,
}: Props) => {
  const router = useRouter();

  const { setLoading } = usePageLoadingStore();
  const { mutate: deleteBookmark, isPending: deletePending } =
    useDeleteFavorite(markerId);
  const { moveLocation } = useMapControl();

  const alertRef = useRef<HTMLButtonElement>(null);

  return (
    <li
      className={`flex w-full items-center p-4 rounded-sm mb-2 duration-100 hover:bg-zinc-700 cursor-pointer hover:scale-95`}
      onClick={() => {
        setLoading(true);
        moveLocation({
          lat: lat as number,
          lng: lng as number,
          isfilter: true,
          markerId: markerId,
        });
        router.push(`/pullup/${markerId}`);
      }}
    >
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
                <BookmarkIcon size={20} isActive />
              )}
            </div>
          </TooltipTrigger>
          <TooltipContent>
            <p>북마크 취소</p>
          </TooltipContent>
        </Tooltip>
      </TooltipProvider>

      <AlertDialog>
        <AlertDialogTrigger asChild>
          <button
            className="hidden"
            ref={alertRef}
            onClick={(e) => {
              e.stopPropagation();
            }}
          >
            북마크 취소
          </button>
        </AlertDialogTrigger>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>정말 취소하시겠습니까?</AlertDialogTitle>
            <AlertDialogDescription className="text-red">
              나중에 다시 등록하실 수 있습니다.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>취소</AlertDialogCancel>
            <AlertDialogAction
              onClick={(e) => {
                e.stopPropagation();
                deleteBookmark();
              }}
            >
              확인
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
          <TooltipTrigger>
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
