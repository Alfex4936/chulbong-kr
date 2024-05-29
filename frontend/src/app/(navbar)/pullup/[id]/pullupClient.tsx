"use client";

import { FacilitiesRes } from "@/api/markers/getFacilities";
import EmojiHoverButton from "@/components/atom/EmojiHoverButton";
import ErrorMessage from "@/components/atom/ErrorMessage";
import LoadingSpinner from "@/components/atom/LoadingSpinner";
import ShareModal from "@/components/common/ShareModal";
import BookmarkIcon from "@/components/icons/BookmarkIcon";
import ChatBubbleIcon from "@/components/icons/ChatBubbleIcon";
import DeleteIcon from "@/components/icons/DeleteIcon";
import DislikeIcon from "@/components/icons/DislikeIcon";
import EditIcon from "@/components/icons/EditIcon";
import RoadViewIcon from "@/components/icons/RoadViewIcon";
import ShareIcon from "@/components/icons/ShareIcon";
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
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import { Skeleton } from "@/components/ui/skeleton";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Textarea } from "@/components/ui/textarea";
import { useToast } from "@/components/ui/use-toast";
import { MOBILE_WIDTH } from "@/constants";
import useInput from "@/hooks/common/useInput";
import useCreateComment from "@/hooks/mutation/comments/useCreateComment";
import useDeleteFavorite from "@/hooks/mutation/favorites/useDeleteFavorite";
import useSetFavorite from "@/hooks/mutation/favorites/useSetFavorite";
import useDeleteMarker from "@/hooks/mutation/marker/useDeleteMarker";
import useMarkerDislike from "@/hooks/mutation/marker/useMarkerDislike";
import useUndoMarkerDislike from "@/hooks/mutation/marker/useUndoMarkerDislike";
import useUpdateDescription from "@/hooks/mutation/marker/useUpdateDescription";
import useFacilitiesData from "@/hooks/query/marker/useFacilitiesData";
import useMarkerData from "@/hooks/query/marker/useMarkerData";
import useWeatherData from "@/hooks/query/marker/useWeatherData";
import useMapStatusStore from "@/store/useMapStatusStore";
import useMobileMapOpenStore from "@/store/useMobileMapOpenStore";
import usePageLoadingStore from "@/store/usePageLoadingStore";
import useRoadviewStatusStore from "@/store/useRoadviewStatusStore";
import useSelectedMarkerStore from "@/store/useSelectedMarkerStore";
import formatDate from "@/utils/formatDate";
import formatFacilities from "@/utils/formatFacilities";
import { useRouter } from "next/navigation";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import IconButton from "./_components/IconButton";
import ImageList from "./_components/ImageList";
import ReviewList from "./_components/ReviewList";

interface Props {
  markerId: number;
}

const PullupClient = ({ markerId }: Props) => {
  const router = useRouter();

  const commentInput = useInput("");
  const [commentError, setCommentError] = useState("");

  const { setLoading } = usePageLoadingStore();

  const alertRef = useRef<HTMLButtonElement>(null);

  const { toast } = useToast();
  const { open: openMobileMap } = useMobileMapOpenStore();
  const { setMarker } = useSelectedMarkerStore();
  const { setPosition } = useMapStatusStore();

  const { open: roadviewOpen, setPosition: setRoadview } =
    useRoadviewStatusStore();

  const { data: marker, isError } = useMarkerData(markerId);
  const { data: facilities } = useFacilitiesData(markerId);

  const { mutateAsync: createComment } = useCreateComment({
    markerId: markerId,
    commentText: commentInput.value,
  });

  const { mutate: dislike, isPending: dislikePending } =
    useMarkerDislike(markerId);
  const { mutate: undoDislike, isPending: undoDislikePending } =
    useUndoMarkerDislike(markerId);

  const { mutate: setFavorite, isPending: setFavoritePending } =
    useSetFavorite(markerId);
  const { mutate: deleteFavorite, isPending: deleteFavoritePending } =
    useDeleteFavorite(markerId);
  const { data: weather, isLoading: weatherLoading } = useWeatherData(
    marker?.latitude as number,
    marker?.longitude as number,
    !!marker
  );

  const { mutate: deleteMarker, isPending: deletePending } = useDeleteMarker({
    id: markerId,
    isRouting: true,
  });

  const [isShare, setIsShare] = useState(false);

  const shareRef = useRef<HTMLDivElement>(null);

  const [isEditDesc, setIsEditDesc] = useState(false);
  const updateDescInput = useInput(marker?.description || "");

  const { mutate: updateDesc } = useUpdateDescription(
    updateDescInput.value,
    markerId
  );

  const changeRoadviewlocation = useCallback(async () => {
    setRoadview(marker?.latitude as number, marker?.longitude as number);
  }, [marker]);

  const facilitiesData = useMemo(() => {
    return formatFacilities(facilities as FacilitiesRes[]);
  }, [facilities]);

  useEffect(() => {
    if (!marker) return;
    setMarker({
      markerId: markerId,
      lat: marker.latitude,
      lng: marker.longitude,
    });
    setPosition(marker.latitude, marker.longitude);

    return () => {
      setMarker(null);
    };
  }, [marker]);

  if (isError)
    return (
      <ErrorMessage
        text="존재하지 않는 위치입니다."
        className="text-lg text-center mt-2"
      />
    );
  if (!marker) return;

  return (
    <div className="mo:mb-12">
      {/* 이미지 배경 */}
      <div
        className="relative w-full h-64 bg-cover bg-center"
        style={{
          backgroundImage: marker.photos
            ? `url(${marker.photos[0].photoUrl})`
            : "url('/metaimg.webp')",
        }}
      >
        {weatherLoading ? (
          <Skeleton className="absolute top-1 left-1 w-28 h-12 bg-black-light-2" />
        ) : (
          <div className="absolute top-1 left-1 flex  items-center py-1 px-2 rounded-sm z-20 bg-black-light-2">
            <img
              className="mr-2"
              src={weather?.iconImage}
              alt={weather?.desc}
            />
            <span className="text-lg font-bold">{weather?.temperature}℃</span>
          </div>
        )}

        <IconButton
          right={10}
          top={10}
          icon={
            setFavoritePending || deleteFavoritePending ? (
              <LoadingSpinner size="xs" />
            ) : (
              <BookmarkIcon isActive={marker.favorited} />
            )
          }
          onClick={() => {
            if (marker.favorited) deleteFavorite();
            else setFavorite();
          }}
          disabled={setFavoritePending || deleteFavoritePending}
        />
        <div ref={shareRef} className="relative">
          <IconButton
            right={10}
            top={50}
            icon={<ShareIcon />}
            onClick={() => setIsShare(true)}
          />
          {isShare && (
            <ShareModal
              link={`${process.env.NEXT_PUBLIC_URL}/pullup/${marker.markerId}`}
              className="absolute top-[90px] right-[10px] z-[200]"
              closeModal={() => setIsShare(false)}
              buttonRef={shareRef}
              lat={marker.latitude}
              lng={marker.longitude}
              filename={marker.address || String(marker.markerId)}
            />
          )}
        </div>
        <IconButton
          right={10}
          top={90}
          icon={<DislikeIcon isActive={marker.disliked || false} />}
          numberState={marker.dislikeCount || 0}
          disabled={dislikePending || undoDislikePending}
          onClick={() => {
            if (marker.disliked) undoDislike();
            else dislike();
          }}
        />
        <IconButton
          right={10}
          top={130}
          icon={<ChatBubbleIcon size={24} selected={false} />}
          onClick={() => router.push(`/pullup/${markerId}/chat`)}
        />
        {marker.isChulbong && (
          <IconButton
            right={10}
            top={170}
            icon={deletePending ? <LoadingSpinner size="xs" /> : <DeleteIcon />}
            onClick={(e) => {
              e.stopPropagation();

              if (!alertRef) return;
              alertRef.current?.click();
            }}
            disabled={deletePending}
          />
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

        <div className="absolute top-0 left-0 w-full h-full bg-black-tp-dark z-10" />
      </div>
      {/* 기구 숫자 카드 */}
      <div className="relative z-30 px-9 -translate-y-14 mo:px-4">
        <div className="h-28">
          <div
            className="bg-black-light-2 flex flex-col justify-center mx-auto 
                        h-full shadow-md w-2/3 py-4 px-10 rounded-sm mo:text-sm mo:px-5 mo:py-2"
          >
            <div className="flex justify-between">
              <span>철봉</span>
              <span
                className={`${
                  facilitiesData.철봉 === "개수 정보 없음"
                    ? "text-[10px]"
                    : "text-normal"
                } flex items-center`}
              >
                {facilitiesData.철봉 === "개수 정보 없음"
                  ? "개수 정보 없음"
                  : `${facilitiesData.철봉}개`}
              </span>
            </div>
            <Separator className="my-2 bg-grey-dark-1" />
            <div className="flex justify-between">
              <span>평행봉</span>
              <span
                className={`${
                  facilitiesData.평행봉 === "개수 정보 없음"
                    ? "text-[10px]"
                    : "text-normal"
                } flex items-center`}
              >
                {facilitiesData.평행봉 === "개수 정보 없음"
                  ? "개수 정보 없음"
                  : `${facilitiesData.평행봉}개`}
              </span>
            </div>
          </div>
        </div>
        {/* 정보 */}
        <div className="mt-4">
          <div className="flex items-center mb-[2px]">
            <span className="mr-1 w-3/4">
              <h1 className="whitespace-normal overflow-visible break-words truncate text-xl">
                {marker.address || "제공되는 주소가 없습니다."}
              </h1>
            </span>
            <button
              onClick={async () => {
                if (window.innerWidth <= MOBILE_WIDTH) {
                  openMobileMap();
                }
                await changeRoadviewlocation();
                roadviewOpen();
              }}
            >
              <RoadViewIcon />
            </button>
          </div>

          <div className="text-xs text-gray-400 mb-5">
            <span>{formatDate(marker.createdAt)}</span>
            <span>({formatDate(marker.updatedAt)}업데이트)</span>
            <span className="mx-1">|</span>
            <button
              className="underline"
              onClick={() => {
                setLoading(true);
                router.push(`/pullup/${markerId}/reportlist`);
              }}
            >
              정보 수정 제안
            </button>
          </div>

          <h2 className="flex items-center">
            {isEditDesc ? (
              <div className="w-full">
                <div className="mb-2">
                  <Input
                    value={updateDescInput.value}
                    onChange={updateDescInput.handleChange}
                  />
                </div>
                <div className="flex">
                  <Button
                    className="border-grey border bg-transparent dark:text-grey hover:bg-white-tp-light hover:border-transparent mr-2"
                    size={"sm"}
                    onClick={() => {
                      updateDesc();
                      setIsEditDesc(false);
                    }}
                  >
                    수정
                  </Button>
                  <Button
                    className="border-grey border bg-transparent dark:text-grey hover:bg-white-tp-light hover:border-transparent"
                    size={"sm"}
                    onClick={() => setIsEditDesc(false)}
                  >
                    취소
                  </Button>
                </div>
              </div>
            ) : (
              <span className="mr-4 w-5/6 break-words">
                {marker.description || "작성된 설명이 없습니다."}
              </span>
            )}
            {marker.isChulbong && !isEditDesc && (
              <button onClick={() => setIsEditDesc(true)}>
                <EditIcon size={15} />
              </button>
            )}
          </h2>
        </div>

        <Separator className="my-3 bg-grey-dark-1" />

        <Tabs defaultValue="photo" className="w-full">
          <TabsList className="w-full">
            <TabsTrigger className="w-1/2" value="photo">
              사진
            </TabsTrigger>
            <TabsTrigger className="w-1/2" value="review">
              리뷰
            </TabsTrigger>
          </TabsList>
          <TabsContent value="photo">
            <ImageList photos={marker.photos} />
          </TabsContent>
          <TabsContent value="review">
            <AlertDialog>
              <AlertDialogTrigger asChild>
                <div className="w-[90%] mx-auto">
                  <EmojiHoverButton
                    emoji="✏️"
                    text="리뷰 작성하기"
                    subText="생각을 공유해 주세요!"
                  />
                </div>
              </AlertDialogTrigger>
              <AlertDialogContent>
                <AlertDialogHeader>
                  <AlertDialogTitle>리뷰 작성하기</AlertDialogTitle>
                </AlertDialogHeader>
                <Textarea
                  className="resize-none"
                  value={commentInput.value}
                  onChange={(e) => {
                    commentInput.handleChange(e);
                    setCommentError("");
                  }}
                />
                <ErrorMessage text={commentError} />
                <AlertDialogFooter>
                  <AlertDialogCancel onClick={commentInput.resetValue}>
                    취소
                  </AlertDialogCancel>
                  <AlertDialogAction
                    onClick={async () => {
                      if (commentInput.value.length > 40) {
                        toast({ description: "40자 이내로 작성해주세요." });
                        commentInput.resetValue();
                        return;
                      }
                      await createComment();
                      commentInput.resetValue();
                    }}
                  >
                    등록
                  </AlertDialogAction>
                </AlertDialogFooter>
              </AlertDialogContent>
            </AlertDialog>
            <ReviewList markerId={marker.markerId} />
          </TabsContent>
        </Tabs>
      </div>
    </div>
  );
};

export default PullupClient;
